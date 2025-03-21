package coditra

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/fauu/coditra/util"
)

type HandlerFunc = func(w http.ResponseWriter, r *http.Request)

func getLookupResultHandler(cache *LookupCache) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := map[string]string{}
		for k, v := range r.URL.Query() {
			queryParams[k] = v[0]
		}

		sourceID := chi.URLParam(r, "source")
		input := chi.URLParam(r, "input")
		source := LookupSources[sourceID]
		if source != nil {
			var result any
			requestHash, cachedResult := cache.query(sourceID, queryParams, input)
			if cachedResult != nil {
				result = *cachedResult
				log.Printf("Lookup request for source %s served from cache", source.Name)
			} else {
				transParams, _ := source.TransformParams(queryParams)
				result, _ = source.DoLookup(SourceEnv, input, transParams)
				if result != nil {
					cache.store(requestHash, result)
				}
			}
			jsonResponse(w, http.StatusOK, result)
		} else {
			jsonResponse(w, http.StatusNotFound, nil)
		}
	}
}

func getDocumentNames(w http.ResponseWriter, _ *http.Request) {
	_, files, err := util.ScanRecursive(Cfg.DocumentsDir, []string{})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// Newest first
	sort.Slice(files, func(i, j int) bool {
		return files[i].Info.ModTime().After(files[j].Info.ModTime())
	})

	var documentNames []string
	for _, file := range files {
		ext := filepath.Ext(file.Path)
		if ext == DocumentExt {
			relativePath := strings.TrimPrefix(file.Path, fmt.Sprintf("%s%c", Cfg.DocumentsDir, os.PathSeparator))
			documentName := strings.TrimSuffix(relativePath, DocumentExt)
			documentNames = append(documentNames, documentName)
		}
	}

	jsonResponse(w, http.StatusOK, documentNames)
}

var htmlBodyRegexp = regexp.MustCompile(`(?s)<body.*?>(.*)</body>`)

func getDocumentByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	filePath := path.Join(Cfg.DocumentsDir, name) + DocumentExt
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err)
		return
	}

	matches := htmlBodyRegexp.FindStringSubmatch(string(fileContent))
	if len(matches) <= 1 {
		errorResponse(w, http.StatusInternalServerError, errors.New("malformed file"))
		return
	}

	jsonResponse(w, http.StatusOK, struct {
		Content string `json:"content"`
	}{matches[1]})
}

func getSetups(w http.ResponseWriter, _ *http.Request) {
	jsonResponse(w, http.StatusOK, Cfg.Setups)
}
