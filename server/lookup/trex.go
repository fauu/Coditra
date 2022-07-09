package lookup

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/fauu/coditra/util"
)

type TrexResult struct {
	Translations []TrexTranslation `json:"translations"`
	Examples     []TrexExample     `json:"examples"`
	Phrases      []TrexPhrase      `json:"phrases"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type TrexTranslation struct {
	Content string `json:"content"`
	URL     string `json:"url"`
}

type TrexExample struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type TrexPhrase struct {
	Content string `json:"content"`
	URL     string `json:"url"`
}

type TrexLookupParams struct {
	SourceLang string
	TargetLang string
}

var trexLangCodeConv = map[string]string{
	"de": "german",
	"en": "english",
	"fr": "french",
	"it": "italian",
	"pl": "polish",
}

var trexBaseURL = "https://tr-ex.me"

func TrexLookupSource() *Source {
	return &Source{
		Name:           "TREX",
		ConstantParams: LookupParams{},
		TransformParams: func(rawParams map[string]string) (any, error) {
			sourceLang := ""
			targetLang := ""
			for k, v := range rawParams {
				switch k {
				case "sourceLang":
					sourceLang = trexLangCodeConv[v]
				case "targetLang":
					targetLang = trexLangCodeConv[v]
				}
			}
			if sourceLang == "" || targetLang == "" {
				return nil, fmt.Errorf("incorrect required lookup params")
			}
			return TrexLookupParams{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			}, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			var sourceLang, targetLang string
			switch v := params.(type) {
			case TrexLookupParams:
				sourceLang = v.SourceLang
				targetLang = v.TargetLang
			default:
				return nil, fmt.Errorf("unpacking lookup params: %v", v)
			}

			urlTemplate := "%s/translation/%s-%s/%s"
			url := fmt.Sprintf(urlTemplate, trexBaseURL, sourceLang, targetLang, input)
			res, err := util.HTTPGet(url, map[string]string{"User-Agent": env.UserAgent})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return nil, err
			}

			var translations []TrexTranslation
			doc.Find(".translations .translation").Each(func(_ int, translationNode *goquery.Selection) {
				translationContent := strings.TrimSpace(translationNode.Text())
				translationURL := fmt.Sprintf("%s?t=%s", url, translationContent)
				translations = append(translations, TrexTranslation{
					Content: translationContent,
					URL:     translationURL,
				})
			})

			var examples []TrexExample
			doc.Find("#contexts .ctx").Each(func(_ int, exampleNode *goquery.Selection) {
				srcHTML, _ := exampleNode.Find(".stc-s .m").Html()
				trgHTML, _ := exampleNode.Find(".stc-t .m").Html()
				examples = append(examples, TrexExample{
					Source: TrexCleanupExampleHTML(srcHTML),
					Target: TrexCleanupExampleHTML(trgHTML),
				})
			})

			var phrases []TrexPhrase
			doc.Find(".context-example").Each(func(_ int, phraseNode *goquery.Selection) {
				phraseURL, _ := phraseNode.Attr("href")
				phrases = append(phrases, TrexPhrase{
					Content: strings.TrimSpace(phraseNode.Text()),
					URL:     trexBaseURL + phraseURL,
				})
			})

			return &TrexResult{
				Translations: translations,
				Examples:     examples,
				Phrases:      phrases,
				SourceURL:    url,
				IsEmpty:      len(examples) == 0,
			}, nil
		},
	}
}

func TrexCleanupExampleHTML(exampleHTML string) string {
	return strings.TrimSpace(html.UnescapeString(AllowFormattedHTMLPolicy.Sanitize(exampleHTML)))
}
