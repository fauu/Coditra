package lookup

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"

	"github.com/fauu/coditra/util"
)

type WrLookupParams struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}

type WrResult struct {
	MainResults   []WrResultEntry `json:"mainResults"`
	CompoundForms []WrResultEntry `json:"compoundForms"`
	OtherDicts    []string        `json:"otherDicts"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type WrResultEntry struct {
	SourceForm   string          `json:"sourceForm"`
	Sense        *WrSense        `json:"sense"`
	Translations []WrTranslation `json:"translations"`
}

type WrSense struct {
	Note  string `json:"note"`
	Sense string `json:"sense"`
}

type WrTranslation struct {
	Note        string `json:"note"`
	Translation string `json:"translation"`
}

var wrOtherDictsPolicy *bluemonday.Policy

func WrLookupSource() *Source {
	return &Source{
		Name:           "WordRef",
		ConstantParams: LookupParams{},
		TransformParams: func(rawParams map[string]string) (any, error) {
			sourceLang := ""
			targetLang := ""
			for k, v := range rawParams {
				switch k {
				case "sourceLang":
					sourceLang = v
				case "targetLang":
					targetLang = v
				}
			}
			if sourceLang == "" || targetLang == "" {
				return nil, fmt.Errorf("incorrect lookup params")
			}
			return WrLookupParams{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			}, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			if wrOtherDictsPolicy == nil {
				wrOtherDictsPolicy = bluemonday.NewPolicy().AllowElements("strong", "em", "div", "span", "ul", "li")
				wrOtherDictsPolicy.AllowAttrs("class").Globally()
			}

			var sourceLang, targetLang string
			switch v := params.(type) {
			case WrLookupParams:
				sourceLang = v.SourceLang
				targetLang = v.TargetLang
			default:
				return nil, fmt.Errorf("unpacking lookup params: %v", v)
			}

			langCombinationString := fmt.Sprintf("%s%s", sourceLang, targetLang)
			urlTemplate := "https://www.wordreference.com/%s/%s"

			url := fmt.Sprintf(urlTemplate, langCombinationString, input)
			res, err := util.HTTPGet(url, map[string]string{"User-Agent": env.UserAgent})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return nil, err
			}

			var result WrResult

			trRowsSel := fmt.Sprintf("tr[id^='%s']", langCombinationString)
			doc.Find(trRowsSel).Each(func(_ int, rowSel *goquery.Selection) {
				var entry WrResultEntry

				sourceForm, sense, translation := wrParseTranslationRow(rowSel)
				entry.SourceForm = sourceForm
				entry.Sense = sense
				entry.Translations = append(entry.Translations, translation)

				for cur := rowSel.Next(); ; cur = cur.Next() {
					if cur.Length() == 0 {
						break
					}
					if cur.HasClass("wrtopsection") {
						break
					}
					if _, hasID := cur.Attr("id"); hasID {
						break
					}
					if cur.Find(".ToEx, .FrEx").Length() > 0 {
						continue
					}
					_, _, translation := wrParseTranslationRow(cur)
					entry.Translations = append(entry.Translations, translation)
				}

				if _, sourceTableHasID := rowSel.Parent().Parent().Attr("id"); sourceTableHasID {
					result.CompoundForms = append(result.CompoundForms, entry)
				} else {
					result.MainResults = append(result.MainResults, entry)
				}
			})

			doc.Find("#otherDicts .entry").Each(func(_ int, entrySel *goquery.Selection) {
				html, err := entrySel.Html()
				if err == nil {
					result.OtherDicts = append(result.OtherDicts, wrOtherDictsPolicy.Sanitize(html))
				}
			})

			result.SourceURL = url
			if len(result.MainResults) == 0 && len(result.CompoundForms) == 0 && len(result.OtherDicts) == 0 {
				result.IsEmpty = true
			}

			return result, nil
		},
	}
}

func wrParseTranslationRow(s *goquery.Selection) (string, *WrSense, WrTranslation) {
	cur := s.Children().First()

	var sourceForm string
	if cur.HasClass("FrWrd") {
		sourceForm = cur.Find("strong").Text()
	}

	cur = cur.Next()

	var sense *WrSense
	senseText := util.GoqueryImmediateText(cur, "")
	if senseText != "" {
		senseText = senseText[1 : len(senseText)-1]
	}
	noteText := cur.Find(".Fr2").Text()
	if senseText != "" || noteText != "" {
		sense = &WrSense{
			Note:  noteText,
			Sense: senseText,
		}
	}

	var translation WrTranslation
	translation.Note = cur.Find(".dsense i").Text()

	cur = cur.Next()

	translation.Translation = util.GoqueryImmediateText(cur, "")

	return sourceForm, sense, translation
}
