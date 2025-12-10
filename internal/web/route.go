package web

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/VineLink-Lab/i18n/internal/parser"
	"golang.org/x/text/language"
)

//go:embed web_file/*
var webFS embed.FS

type RouteHandler struct {
	parser *parser.Parser
}

func NewRouteHandler(parser *parser.Parser) *RouteHandler {
	return &RouteHandler{
		parser: parser,
	}
}

func (rh *RouteHandler) RegisterRoutes() {
	webFS, _ := fs.Sub(webFS, "web_file")
	http.Handle("/", http.FileServer(http.FS(webFS)))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	rh.registerBundlesRoutes()
}

func (rh *RouteHandler) StartServer(address string) error {
	return http.ListenAndServe(address, nil)
}

func writeJSONResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonData)
}

func parseJSONBody(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func parseLanguageTag(langStr string) (language.Tag, error) {
	return language.Parse(langStr)
}
