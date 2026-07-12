package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/raaj2493/Shortly/Backend/internals/services"
)

// URLHandler handles HTTP requests for URLs.
type URLHandler struct {
	urlService *services.URLService
}

// NewURLHandler creates a new handler.
func NewURLHandler(urlService *services.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

// CreateShortURL handles POST /api/urls.
// Expects JSON: {"original_url": "https://example.com/long"}
// Returns: {"short_code": "abc123", "short_url": "http://localhost:8080/abc123"}
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Only accept POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var req struct {
		OriginalURL string `json:"original_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.OriginalURL == "" {
		http.Error(w, "original_url is required", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call service
	url, err := h.urlService.CreateShortURL(req.OriginalURL, userID)
	if err != nil {
		http.Error(w, "Failed to create short URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response
	shortURL := "http://localhost:8080/" + url.ShortCode // In production, use config base URL
	resp := map[string]string{
		"short_code": url.ShortCode,
		"short_url":  shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Redirect handles GET /{shortCode}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Extract short code from URL path
	shortCode := r.PathValue("shortCode") // Requires Go 1.22+ path parameters
	if shortCode == "" {
		http.Error(w, "Missing short code", http.StatusBadRequest)
		return
	}

	// Look up original URL
	originalURL, err := h.urlService.GetOriginalURL(shortCode)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect with 301 (permanent) to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}