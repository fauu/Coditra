package lookup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/fauu/coditra/util"
)

type GarzantiResult struct {
	Word            string `json:"word"`
	Hyphenation     string `json:"hyphenation"`
	GrammaticalType string `json:"grammaticalType"`
	EtymologyHTML   string `json:"etymologyHtml"`
	DefinitionHTML  string `json:"definitionHtml"`

	SourceURL string `json:"sourceUrl"`
	IsEmpty   bool   `json:"isEmpty"`
}

type GarzantiLemmaResponse struct {
	Word             string `json:"word"`
	Sillabazione     string `json:"sillabazione"`
	TipoGrammaticale string `json:"tipogrammaticale"`
	Polirematica     string `json:"polirematica"`
	Etimologia       string `json:"etimologia"`
	Testo            string `json:"testo"`
}

var GarzantiSessionCookie *http.Cookie

func GarzantiLookupSource() *Source {
	return &Source{
		Name:           "Garzanti",
		ConstantParams: LookupParams{SourceLang: &it},
		TransformParams: func(rawParams map[string]string) (any, error) {
			return nil, nil
		},
		DoLookup: func(env SourceEnv, input string, params any) (any, error) {
			baseURL := "https://www.garzantilinguistica.it"

			if GarzantiSessionCookie == nil {
				res, err := util.HTTPGet(baseURL, map[string]string{"User-Agent": env.UserAgent})
				if err != nil {
					return nil, err
				}
				defer res.Body.Close()
				cookies := res.Cookies()
				if len(cookies) == 0 {
					return nil, errors.New("no Garzanti session cookie returned")
				}
				GarzantiSessionCookie = cookies[0]
				GarzantiSessionCookie.Path = "" // So that it serializes properly when sent back
			}

			url := fmt.Sprintf("%s/ricerca/", baseURL)
			res, err := util.HTTPPost(url,
				map[string]string{
					"User-Agent": env.UserAgent,
					// Without the cookie it mixes in results in other languages.
					// Should be using Request.AddCookie() really.
					"Cookie": GarzantiSessionCookie.String(),
				},
				map[string]string{"form_search_word_id": "", "form_search_word": input})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return nil, err
			}

			currentLemmaSel := doc.Find(".current.lemma").First()
			lemmaID, _ := currentLemmaSel.Attr("rel")

			url = fmt.Sprintf("%s/wp-content/themes/garzantilinguistica_new/services/svr-search-lemma.php", baseURL)
			res, err = util.HTTPPost(url,
				map[string]string{
					"User-Agent": env.UserAgent,
					"Cookie":     GarzantiSessionCookie.String(),
				},
				map[string]string{"idlemma": lemmaID})
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			bodyBytes = bytes.TrimPrefix(bodyBytes, []byte("\xef\xbb\xbf")) // BOM
			bodyStr := string(bodyBytes)
			var lemmaResponse GarzantiLemmaResponse
			err = json.Unmarshal([]byte(bodyStr), &lemmaResponse)
			if err != nil {
				return nil, err
			}

			defHTML := lemmaResponse.Testo
			if len(defHTML) > 0 {
				defHTML = fmt.Sprintf("%s</p>", defHTML[4:])
				defHTML = strings.ReplaceAll(defHTML, "<p></p>", "")
				defHTML = strings.ReplaceAll(defHTML, " |", " | ")
			}

			return GarzantiResult{
				Word:            lemmaResponse.Word,
				Hyphenation:     lemmaResponse.Sillabazione,
				GrammaticalType: lemmaResponse.TipoGrammaticale,
				EtymologyHTML:   lemmaResponse.Etimologia,
				DefinitionHTML:  defHTML,

				SourceURL: fmt.Sprintf("%s/ricerca/?q=%s", baseURL, StripTagsPolicy.Sanitize(lemmaResponse.Word)),
				IsEmpty:   len(defHTML) == 0,
			}, nil
		},
	}
}
