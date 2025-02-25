package custommiddleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

// ContextKey is used to store values in context
type ContextKey string

const UserIDKey ContextKey = "user_id"

// AuthMiddleware verifies Firebase ID tokens and extracts the user ID
func AuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, `{"error": "Missing authentication token"}`, http.StatusUnauthorized)
				return
			}

			// Remove 'Bearer ' prefix if present
			tokenString = strings.TrimSpace(strings.Replace(tokenString, "Bearer ", "", 1))

			// Verify Firebase token
			token, err := authClient.VerifyIDToken(r.Context(), tokenString)
			if err != nil {
				http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Extract user ID
			uid, ok := token.Claims["user_id"].(string)
			if !ok {
				http.Error(w, `{"error": "User ID missing in token"}`, http.StatusUnauthorized)
				return
			}

			// Store user ID in request context
			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
