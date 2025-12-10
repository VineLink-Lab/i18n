package web

import (
	"net/http"

	"github.com/VineLink-Lab/i18n/internal/parser"
)

func (rh *RouteHandler) registerBundlesRoutes() {
	http.HandleFunc("/api/bundles", rh.handleBundles)
	http.HandleFunc("/api/bundles/content", rh.handleBundleContent)
	http.HandleFunc("/api/languages", rh.handleLanguages)
	http.HandleFunc("/api/message/update", rh.handleUpdateMessage)
	http.HandleFunc("/api/message/delete", rh.handleDeleteMessage)
	http.HandleFunc("/api/message/add", rh.handleAddMessage)
}

func (rh *RouteHandler) handleBundles(w http.ResponseWriter, r *http.Request) {
	bundles := rh.parser.GetAvailableBundles()
	writeJSONResponse(w, bundles)
}

func (rh *RouteHandler) handleLanguages(w http.ResponseWriter, r *http.Request) {
	languages := rh.parser.GetSupportedLanguages()
	langStrings := make([]string, len(languages))
	for i, lang := range languages {
		langStrings[i] = lang.String()
	}
	writeJSONResponse(w, langStrings)
}

func (rh *RouteHandler) handleBundleContent(w http.ResponseWriter, r *http.Request) {
	bundleName := r.URL.Query().Get("bundle")
	if bundleName == "" {
		bundleName = "."
	}

	content, err := rh.parser.GetBundleContent(bundleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to frontend-friendly format
	response := make(map[string]map[string]interface{})
	for lang, messages := range content {
		langStr := lang.String()
		response[langStr] = make(map[string]interface{})
		for key, message := range messages {
			response[langStr][key] = message
		}
	}

	writeJSONResponse(w, response)
}

// UpdateMessageRequest is the request to update a message
type UpdateMessageRequest struct {
	Bundle      string `json:"bundle"`
	Key         string `json:"key"`
	Lang        string `json:"lang"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

// DeleteMessageRequest is the request to delete a message
type DeleteMessageRequest struct {
	Bundle string `json:"bundle"`
	Key    string `json:"key"`
}

// AddMessageRequest is the request to add a new message
type AddMessageRequest struct {
	Bundle       string            `json:"bundle"`
	Key          string            `json:"key"`
	Translations map[string]string `json:"translations"` // lang -> content
	Description  string            `json:"description"`
}

func (rh *RouteHandler) handleUpdateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateMessageRequest
	if err := parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse language tag
	lang, err := parseLanguageTag(req.Lang)
	if err != nil {
		http.Error(w, "Invalid language tag", http.StatusBadRequest)
		return
	}

	// Build message
	message := parser.Message{
		Content:     req.Content,
		Description: req.Description,
	}

	// Update message
	err = rh.parser.UpdateMessage(req.Bundle, req.Key, lang, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]string{"status": "ok"})
}

func (rh *RouteHandler) handleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteMessageRequest
	if err := parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Delete message
	err := rh.parser.DeleteMessage(req.Bundle, req.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]string{"status": "ok"})
}

func (rh *RouteHandler) handleAddMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddMessageRequest
	if err := parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add translation for each language
	for langStr, content := range req.Translations {
		lang, err := parseLanguageTag(langStr)
		if err != nil {
			http.Error(w, "Invalid language tag: "+langStr, http.StatusBadRequest)
			return
		}

		message := parser.Message{
			Content:     content,
			Description: req.Description,
		}

		err = rh.parser.UpdateMessage(req.Bundle, req.Key, lang, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSONResponse(w, map[string]string{"status": "ok"})
}
