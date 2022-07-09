package lookup

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"github.com/fauu/coditra/util"
)

type SynonimyPlResult struct {
	Title                 string                `json:"title"`
	SynonymGroups         [][]SynonimyPlSynonym `json:"synonymGroups"`
	SuggestedAlternatives []string              `json:"suggestedAlternatives"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type SynonimyPlSynonym struct {
	Synonym string `json:"synonym"`
	Extra   string `json:"extra"`
}

var synonimyplTitleRegexp = regexp.MustCompile(`"(.*?)"`)

func SynonimyPlLookupSource() *Source {
	return &Source{
		Name:           "Synonimy.pl",
		ConstantParams: LookupParams{SourceLang: &pl},
		TransformParams: func(rawParams map[string]string) (any, error) {
			return nil, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			urlTemplate := "https://www.synonimy.pl/synonim/%s/"

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

			var result SynonimyPlResult

			titleMatch := synonimyplTitleRegexp.FindStringSubmatch(doc.Find("title").Text())
			if len(titleMatch) > 1 {
				result.Title = titleMatch[1]
			}

			termSel := doc.Find(".term")

			if termSel.Length() == 0 {
				suggested := doc.Find(".load_word").Map(func(_ int, wordSel *goquery.Selection) string {
					return wordSel.Text()
				})
				return &SynonimyPlResult{SuggestedAlternatives: suggested}, nil
			}

			termSel.Find("dd > span").Each(func(_ int, groupSel *goquery.Selection) {
				var group []SynonimyPlSynonym
				for cur := groupSel; ; {
					synonym := SynonimyPlSynonym{
						Synonym: cur.Children().Get(1).FirstChild.Data,
						Extra:   cur.Children().First().Text(),
					}
					group = append(group, synonym)
					if cur.Children().Length() != 4 {
						break
					}
					cur = cur.Children().Last()
					if _, hasClass := cur.Attr("class"); hasClass { // Ad container
						break
					}
				}
				result.SynonymGroups = append(result.SynonymGroups, group)
			})

			result.SourceURL = url
			result.IsEmpty = false

			return result, nil
		},
	}
}
