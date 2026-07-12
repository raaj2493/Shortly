package middleware

import (
	"encoding/json"
	"net/http"
)

// URLHandler holds dependencies for URL endpoints.
type URLHandler struct {
	// We'll add a *service.URLService later.
}

// NewURLHandler creates a new URLHandler (for now without a service).
func NewURLHandler() *URLHandler {
	return &URLHandler{}
}

// CreateShortURL handles POST /api/urls (protected). Placeholder returns a dummy response.
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	// Placeholder: we'll implement in Module 5
	writeJSON(w, http.StatusOK, map[string]string{"message": "CreateShortURL not implemented yet"})
}

// RedirectShortURL handles GET /{shortCode}. Placeholder.
func (h *URLHandler) RedirectShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	// Placeholder: we'll implement in Module 5
	writeJSON(w, http.StatusOK, map[string]string{"message": "RedirectShortURL not implemented yet"})
}

// writeJSON is a helper that sets Content-Type and status code, then encodes data as JSON.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}