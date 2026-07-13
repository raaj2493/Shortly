package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"github.com/raaj2493/Shortly/Backend/internals/services"

)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

// CreateShortURL handles POST /api/urls
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the incoming JSON
	var request struct {
		OriginalURL string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if request.OriginalURL == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "original_url is required",
		})
		return
	}

	// 2. Call the service to generate a short code
	shortCode, err := h.service.Shorten(r.Context(), request.OriginalURL)
	if err != nil {
		log.Printf("Error shortening URL: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create short URL",
		})
		return
	}

	// 3. Return success with the short code
	writeJSON(w, http.StatusCreated, map[string]string{
		"short_code": shortCode,
	})
}

// RedirectToOriginal handles GET /{shortCode}
func (h *URLHandler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the short code from the URL path
	shortCode := r.PathValue("shortCode")
	if shortCode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "short code is required",
		})
		return
	}

	// 2. Look up the original URL
	originalURL, err := h.service.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		log.Printf("Error looking up short code %q: %v", shortCode, err)
		if errors.Is(err, errors.New("URL not found")) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "Short URL not found",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// 3. Issue a 301 Moved Permanently redirect
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}



// writeJSON is a small helper to avoid boilerplate
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

