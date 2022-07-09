package lookup

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/fauu/coditra/util"
)

type PwnKorpusResult struct {
	Entries []PwnKorpusResultEntry `json:"entries"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type PwnKorpusResultEntry struct {
	PreFragment  string `json:"preFragment"`
	TheWord      string `json:"theWord"`
	PostFragment string `json:"postFragment"`
	DetailsURL   string `json:"detailsUrl"`
}

func PwnKorpusLookupSource() *Source {
	return &Source{
		Name:           "Korpus PWN",
		ConstantParams: LookupParams{SourceLang: &pl},
		TransformParams: func(rawParams map[string]string) (any, error) {
			return nil, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			urlTemplate := "https://sjp.pwn.pl/korpus/szukaj/%s.html"

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

			var resultEntries []PwnKorpusResultEntry
			doc.Find(".sjp-korpus-lista li.row > div").Each(func(_ int, s *goquery.Selection) {
				partsSel := s.Children()

				preFragment := partsSel.Eq(0).Text()
				postFragment := partsSel.Eq(2).Text()
				preFragment = strings.Replace(preFragment, "...", "…", -1)
				postFragment = strings.Replace(postFragment, "...", "…", -1)

				theWordNode := partsSel.Eq(1)
				theWordAnchorHref, _ := theWordNode.Children().First().Attr("href")

				resultEntries = append(resultEntries, PwnKorpusResultEntry{
					PreFragment:  preFragment,
					TheWord:      theWordNode.Text(),
					DetailsURL:   theWordAnchorHref,
					PostFragment: postFragment,
				})
			})

			result := PwnKorpusResult{
				Entries:   resultEntries,
				SourceURL: url,
				IsEmpty:   len(resultEntries) == 0,
			}

			return result, nil
		},
	}
}
