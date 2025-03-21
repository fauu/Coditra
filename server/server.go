package coditra

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"github.com/fauu/coditra/lookup"
)

//go:embed client-dist
var outsideClientDistFS embed.FS

var clientDistFS fs.FS

func init() {
	clientDistFS, _ = fs.Sub(outsideClientDistFS, "client-dist")
}

var CORSOptions = cors.Options{
	AllowedOrigins: []string{"http://localhost:5000"},
	AllowedMethods: []string{"GET", "OPTIONS", "POST", "PUT", "DELETE"},
	AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding",
		"X-CSRF-Token", "Authorization", "Access-Control-Request-Headers",
		"Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent",
		"Referer", "Cache-Control", "X-header"},
	OptionsPassthrough: true,
}

var LookupSources map[string]*lookup.Source
var Cfg *Config
var SourceEnv lookup.SourceEnv

func RunServer() error {
	LookupSources = map[string]*lookup.Source{
		"pwn":        lookup.PwnLookupSource(),
		"pwnkorpus":  lookup.PwnKorpusLookupSource(),
		"rc":         lookup.RcLookupSource(),
		"synonimypl": lookup.SynonimyPlLookupSource(),
		"wr":         lookup.WrLookupSource(),
		"lingea":     lookup.LingeaLookupSource(),
		"garzanti":   lookup.GarzantiLookupSource(),
		"trex":       lookup.TrexLookupSource(),
	}

	var err error
	Cfg, err = LoadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %v", err)
	}

	SourceEnv = lookup.SourceEnv{
		UserAgent: Cfg.UserAgent,
	}

	lookupCache := NewLookupCache()

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/docs", getDocumentNames)
		r.Get("/docs/{name}", getDocumentByName)
		r.Get("/setups", getSetups)
		r.Get("/{source}/{input}", getLookupResultHandler(&lookupCache))
	})
	FileServer(r, "/", http.FS(clientDistFS))

	fmt.Printf("Starting coditra server at %s\n", Addr)
	err = http.ListenAndServe(Addr, cors.New(CORSOptions).Handler(r))
	if err != nil {
		return fmt.Errorf("running http server: %v", err)
	}

	return nil
}

// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		myfs := http.StripPrefix(pathPrefix, http.FileServer(root))
		myfs.ServeHTTP(w, r)
	})
}
