package lookup

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/fauu/coditra/util"
)

type LingeaLookupParams struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}

type LingeaResult struct {
	Title        string                 `json:"title"`
	FemSuffix    string                 `json:"femSuffix"`
	Related      []string               `json:"related"`
	MorfBlocks   []LingeaMorfBlock      `json:"morfBlocks"`
	Phrases      LingeaMorfBlockSegment `json:"phrases"`
	KeywordTerms []LingeaKeywordTerm    `json:"keywordTerms"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type LingeaMorfBlock struct {
	Morf     string                   `json:"morf"`
	Segments []LingeaMorfBlockSegment `json:"segments"`
}

type LingeaMorfBlockSegment struct {
	Definition  *string           `json:"definition"`
	Expressions []LingeaTwoColumn `json:"expressions"`
	Examples    []LingeaTwoColumn `json:"examples"`
}

type LingeaKeywordTerm struct {
	Keyword   string          `json:"keyword"`
	TwoColumn LingeaTwoColumn `json:"twoColumn"`
}

type LingeaTwoColumn struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

var lingeaLangCodeConv = map[string][]string{
	"de": {"niemiecko", "niemiecki"},
	"en": {"angielsko", "angielski"},
	"fr": {"francusko", "francuski"},
	"it": {"wlosko", "wloski"},
	"pl": {"polsko", "polski"},
}

var lingeaMorfConv = map[string]string{
	"adj":  "przymiotnik",
	"adv":  "przysłówek",
	"art":  "przedimek",
	"conj": "spójnik",
	"f":    "rzeczownik",
	"m":    "rzeczownik",
	"n":    "rzeczownik",
	"prep": "przyimek",
	"pron": "zaimek",
	"v":    "czasownik",
}

func LingeaLookupSource() *Source {
	return &Source{
		Name:           "Lingea",
		ConstantParams: LookupParams{},
		TransformParams: func(rawParams map[string]string) (any, error) {
			sourceLang := ""
			targetLang := ""
			for k, v := range rawParams {
				switch k {
				case "sourceLang":
					sourceLang = lingeaLangCodeConv[v][0]
				case "targetLang":
					targetLang = lingeaLangCodeConv[v][1]
				}
			}
			if sourceLang == "" || targetLang == "" {
				return nil, fmt.Errorf("missing required lookup params")
			}
			return LingeaLookupParams{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			}, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			var sourceLang, targetLang string
			switch v := params.(type) {
			case LingeaLookupParams:
				sourceLang = v.SourceLang
				targetLang = v.TargetLang
			default:
				return nil, fmt.Errorf("unpacking lookup params: %v", v)
			}

			urlTemplate := "https://slowniki.lingea.pl/%s-%s/%s"

			url := fmt.Sprintf(urlTemplate, sourceLang, targetLang, input)
			res, err := util.HTTPGet(url, map[string]string{"User-Agent": env.UserAgent})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return nil, err
			}

			result := LingeaResult{}

			headSel := doc.Find(".head")

			result.Title = strings.ReplaceAll(headSel.Find(".lex_ful_entr").Text(), "*", "")
			result.FemSuffix = strings.TrimPrefix(headSel.Find(".lex_ful_vfem").Text(), ", ")
			headMorf := headSel.Find(".lex_ful_morf").Text()

			segmentsSel := headSel.NextAll()
			morfBlock := LingeaMorfBlock{}
			if headMorf != "" {
				morfBlock.Morf = lingeaParseMorf(headMorf)
			}

			segmentsSel.Each(func(_ int, s *goquery.Selection) {
				if s.Children().Length() == 1 {
					morfSel := s.Find(".lex_ful_morf")
					if morfSel.Length() == 1 {
						if morfBlock.Morf != "" {
							result.MorfBlocks = append(result.MorfBlocks, morfBlock)
							morfBlock = LingeaMorfBlock{}
						}

						morfBlock.Morf = lingeaParseMorf(morfSel.Text())
					}
					return
				}
				morfBlock.Segments = append(morfBlock.Segments, lingeaParseMorfBlockSegment(s))
			})
			if morfBlock.Morf == "phr" {
				result.Phrases = morfBlock.Segments[0]
			} else {
				result.MorfBlocks = append(result.MorfBlocks, morfBlock)
			}

			doc.Find(".lex_ftx_sens").Each(func(_ int, s *goquery.Selection) {
				keyword := s.Find(".lex_ftx_entr").Text()
				sourceSel := s.Find(".lex_ftx_samp2s")
				targetSel := s.Find(".lex_ftx_samp2t")
				source, _ := sourceSel.Html()
				target, _ := targetSel.Html()
				twoColumn := LingeaTwoColumn{Source: source, Target: target}
				result.KeywordTerms = append(result.KeywordTerms, LingeaKeywordTerm{Keyword: keyword, TwoColumn: twoColumn})
			})

			doc.Find(".Find .bspan").Each(func(_ int, s *goquery.Selection) {
				related := s.Text()
				related = strings.ReplaceAll(related, "*", "")
				if related != result.Title {
					result.Related = append(result.Related, related)
				}
			})

			result.SourceURL = url
			if result.Title == "" {
				result.IsEmpty = true
			}

			return &result, nil
		},
	}
}

func lingeaParseMorf(rawMorf string) string {
	rawMorf = strings.TrimSpace(rawMorf)
	if morf, ok := lingeaMorfConv[rawMorf]; ok {
		return morf
	}
	return rawMorf
}

func lingeaParseMorfBlockSegment(rootSel *goquery.Selection) LingeaMorfBlockSegment {
	segment := LingeaMorfBlockSegment{}

	cur := rootSel.Children().First().Next()

	definitionSel := cur.Find(".lex_ful_tran").First()
	if definitionSel.Length() != 0 {
		definition, _ := definitionSel.Html()
		segment.Definition = &definition
	}

	cur.Find(".lex_ful_coll2").Each(func(_ int, s *goquery.Selection) {
		sourceSel := s.Find(".lex_ful_coll2s")
		targetSel := s.Find(".lex_ful_coll2t")
		source, _ := sourceSel.Html()
		target, _ := targetSel.Html()
		expression := LingeaTwoColumn{Source: source, Target: target}
		segment.Expressions = append(segment.Expressions, expression)
	})

	cur.Find(".lex_ful_samp2").Each(func(_ int, s *goquery.Selection) {
		sourceSel := s.Find(".lex_ful_samp2s")
		targetSel := s.Find(".lex_ful_samp2t")
		source, _ := sourceSel.Html()
		target, _ := targetSel.Html()
		example := LingeaTwoColumn{Source: source, Target: target}
		segment.Examples = append(segment.Examples, example)
	})

	return segment
}
