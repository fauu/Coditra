package lookup

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/fauu/coditra/util"
)

type RcResult struct {
	Translations []RcTranslation `json:"translations"`
	Examples     []RcExample     `json:"examples"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type RcTranslation struct {
	Content string `json:"content"`
	URL     string `json:"url"`
}

type RcExample struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type RcLookupParams struct {
	SourceLang string
	TargetLang string
}

var rcLangCodeConv = map[string]string{
	"de": "german",
	"en": "english",
	"fr": "french",
	"it": "italian",
	"pl": "polish",
}

var rcBaseURL = "https://context.reverso.net"

func RcLookupSource() *Source {
	return &Source{
		Name:           "RevCon",
		ConstantParams: LookupParams{},
		TransformParams: func(rawParams map[string]string) (any, error) {
			sourceLang := ""
			targetLang := ""
			for k, v := range rawParams {
				switch k {
				case "sourceLang":
					sourceLang = rcLangCodeConv[v]
				case "targetLang":
					targetLang = rcLangCodeConv[v]
				}
			}
			if sourceLang == "" || targetLang == "" {
				return nil, fmt.Errorf("incorrect required lookup params")
			}
			return RcLookupParams{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			}, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			var sourceLang, targetLang string
			switch v := params.(type) {
			case RcLookupParams:
				sourceLang = v.SourceLang
				targetLang = v.TargetLang
			default:
				return nil, fmt.Errorf("unpacking lookup params: %v", v)
			}

			urlTemplate := "%s/translation/%s-%s/%s"
			url := fmt.Sprintf(urlTemplate, rcBaseURL, sourceLang, targetLang, input)
			res, err := util.HTTPGet(url, map[string]string{"User-Agent": env.UserAgent})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return nil, err
			}

			var translations []RcTranslation
			doc.Find("#translations-content .translation.dict").Each(func(_ int, translationNode *goquery.Selection) {
				translationContent := strings.TrimSpace(translationNode.Text())
				translationURL := fmt.Sprintf("%s#%s", url, translationContent)
				translations = append(translations, RcTranslation{Content: translationContent, URL: translationURL})
			})

			var examples []RcExample
			doc.Find(".example").Each(func(_ int, exampleNode *goquery.Selection) {
				srcHTML, _ := exampleNode.Find(".src").First().Find(".text").First().Html()
				trgHTML, _ := exampleNode.Find(".trg").First().Find(".text").First().Html()
				examples = append(examples, RcExample{Source: RcCleanupExampleHTML(srcHTML), Target: RcCleanupExampleHTML(trgHTML)})
			})

			return &RcResult{Translations: translations, Examples: examples, SourceURL: url, IsEmpty: len(examples) == 0}, nil
		},
	}
}

func RcCleanupExampleHTML(exampleHTML string) string {
	return strings.TrimSpace(html.UnescapeString(AllowFormattedHTMLPolicy.Sanitize(exampleHTML)))
}
