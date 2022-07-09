package lookup

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"github.com/fauu/coditra/util"
)

type PwnResult struct {
	SjpResult         PwnSjpResult         `json:"sjpResult"`
	DoroszewskiResult PwnDoroszewskiResult `json:"doroszewskiResult"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type PwnSjpResult struct {
	Entries []PwnSjpResultEntry `json:"entries"`
}

type PwnSjpResultEntry struct {
	Title       string             `json:"title"`
	Definitions []PwnSjpDefinition `json:"definitions"`
}

type PwnSjpDefinition struct {
	Qualifier *string `json:"qualifier"`
	Content   string  `json:"content"`
	Xref      *string `json:"xref"`
}

type PwnDoroszewskiResult struct {
	Entries []PwnDoroszewskiResultEntry `json:"entries"`
}

type PwnDoroszewskiResultEntry struct {
	Title           string   `json:"title"`
	ImgFragmentURLs []string `json:"imgFragmentUrls"`
}

func PwnLookupSource() *Source {
	return &Source{
		Name:           "SJP PWN",
		ConstantParams: LookupParams{SourceLang: &pl},
		TransformParams: func(rawParams map[string]string) (any, error) {
			return nil, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			urlTemplate := "https://sjp.pwn.pl/szukaj/%s.html"

			url := fmt.Sprintf(urlTemplate, input)
			res, err := util.HTTPGet(url, map[string]string{"User-Agent": env.UserAgent})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return nil, err
			}

			sjpResult := PwnSjpResult{}
			doc.Find(".sjp-wyniki .entry-body .ribbon-element").Each(func(_ int, s *goquery.Selection) {
				title := s.Find(".tytul-grupa").Text()
				if title == "" {
					titleParts := s.Find(".tytul").Map(func(_ int, t *goquery.Selection) string {
						return t.Text()
					})
					title = strings.Join(titleParts, ", ")
				}

				defSels := s.Find(".znacz")
				if defSels.Length() == 0 {
					defSels = s
				}
				var defs []PwnSjpDefinition
				defSels.Each(func(_ int, defSel *goquery.Selection) {
					defs = append(defs, pwnParseDefinition(defSel))
				})

				sjpResult.Entries = append(sjpResult.Entries, PwnSjpResultEntry{
					Title:       title,
					Definitions: defs,
				})
			})

			var wg sync.WaitGroup

			doroEntrySel := doc.Find(".sjp-doroszewski-wyniki .entry-body a.anchor-title")
			doroEntries := make([]PwnDoroszewskiResultEntry, doroEntrySel.Length())
			doroEntrySel.Each(func(i int, s *goquery.Selection) {
				wg.Add(1)
				go func(i int, s *goquery.Selection) {
					defer wg.Done()
					href, _ := s.Attr("href")
					imgFragmentURLs, _ := pwnGetDoroszewskiImgFragmentURLs(env, href)
					doroEntries[i] = PwnDoroszewskiResultEntry{
						Title:           s.Text(),
						ImgFragmentURLs: imgFragmentURLs,
					}
				}(i, s)
			})
			wg.Wait()

			result := PwnResult{
				SjpResult:         sjpResult,
				DoroszewskiResult: PwnDoroszewskiResult{Entries: doroEntries},
				SourceURL:         url,
			}
			result.IsEmpty = len(result.SjpResult.Entries) == 0 && len(result.DoroszewskiResult.Entries) == 0

			return result, nil
		},
	}
}

var pwnDefinitionContentRegexp = regexp.MustCompile("«(.+?)»")

var pwnXrefSuffixRegexp = regexp.MustCompile(` w zn\. \d`)

func pwnParseDefinition(s *goquery.Selection) PwnSjpDefinition {
	var qualifier *string
	var content string
	var xref *string

	maybeXrefIndicatorSel := s.Find("i").First()
	if maybeXrefIndicatorSel.Text() == "zob." {
		xrefAnchorSel := s.Find(".anchor").First()
		if xrefAnchorSel.Length() > 0 {
			content = xrefAnchorSel.Text()
			nonNullXref := pwnXrefSuffixRegexp.ReplaceAllLiteralString(content, "")
			xref = &nonNullXref
		} else {
			// NOTE: This presumably handles the local ref case, but probably leaves content not formatted properly
			fullTextParts := strings.Split(s.Text(), "zob. ")
			content = strings.TrimSpace(fullTextParts[1])
		}
	} else {
		// ASSUMPTION: xrefs can't have qualifiers
		qualifierText := s.Find(".kwal").Text()
		if qualifierText != "" {
			qualifier = &qualifierText
		}

		fullText := s.Text()
		contentMatches := pwnDefinitionContentRegexp.FindStringSubmatch(fullText)
		if len(contentMatches) < 2 {
			content = strings.TrimSpace(fullText)
		} else {
			content = contentMatches[1]
		}
	}

	return PwnSjpDefinition{
		Qualifier: qualifier,
		Content:   content,
		Xref:      xref,
	}
}

func pwnGetDoroszewskiImgFragmentURLs(env SourceEnv, url string) ([]string, error) {
	res, err := util.HTTPGet(url, map[string]string{"User-Agent": env.UserAgent})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc.Find("img.img-dorosz").Map(func(_ int, s *goquery.Selection) string {
		src, _ := s.Attr("src")
		return src
	}), nil
}
