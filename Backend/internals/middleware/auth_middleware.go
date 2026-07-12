package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// UserIDKey is the context key used to store the authenticated user's ID.
// We use a custom type to avoid collisions with other context keys.
type contextKey string
const UserIDKey contextKey = "userID"

// AuthMiddleware returns a middleware that protects routes.
// It expects a JWT in the Authorization header as "Bearer <token>".
// It validates the token, checks that a corresponding session exists in Redis,
// and injects the user ID into the request context.
func AuthMiddleware(jwtSecret string, redisSessionChecker func(userID int) (bool, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"missing or malformed token"}`, http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// 2. Parse and validate the JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// 3. Extract claims (we expect the "sub" claim to be the user ID)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, `{"error":"invalid subject claim"}`, http.StatusUnauthorized)
				return
			}

			// The user ID was stored as a string, so we need to convert it to an int.
			var userID int
			_, err = fmt.Sscanf(sub, "%d", &userID)
			if err != nil {
				http.Error(w, `{"error":"invalid user ID in token"}`, http.StatusUnauthorized)
				return
			}

			// 4. Verify the session exists in Redis
			// redisSessionChecker is a function that checks for session:<userID> existence.
			exists, err := redisSessionChecker(userID)
			if err != nil || !exists {
				http.Error(w, `{"error":"session not found"}`, http.StatusUnauthorized)
				return
			}

			// 5. Inject user ID into context and proceed to the next handler
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}