package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/raaj2493/Shortly/Backend/internals/services"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func(h *AuthHandler) Register (w http.ResponseWriter , r *http.Request) {
	// only accept post route 
	if r.Method != http.MethodPost {
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// decode the request to a user struct
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// 3. Basic validation
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
		return
	}

	// 4. Call the service layer
	user, err := h.AuthService.Register(r.Context() , req.Email, req.Password )
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "registration failed"})
		return
	}

	// 5. Return success response
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user_id": user.ID,
	})

}


func (h *AuthHandler) Login (w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
		return
	}


	// Call service to login, which returns a JWT token string
	token, err := h.AuthService.Login(r.Context(),req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}

	// Return the token
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}


// Logout handles POST /logout. This is a protected route.



// writeJSON is a small helper that sets Content-Type and status code, then encodes data as JSON.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}