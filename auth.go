package main

import (
	"context"
	"kithli-api/firebase"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

// Auth checks the request's JWT token for validity.
// Returns status 401 for invalid JWT tokens, or continues to the next handler.
func Auth(firebaseClient *firebase.FirebaseClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the Firebase Auth client from the global initialization
			firebaseAuth := firebaseClient.FirebaseAuth
			if firebaseAuth == nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, "Firebase Auth client not initialized")
				return
			}

			// Extract and verify the ID token from the Authorization header
			header := r.Header.Get("Authorization")
			idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
			if idToken == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, "Missing or invalid Authorization header")
				return
			}

			token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": err.Error()})
				return
			}

			// Pass the verified token in the request context
			ctx := context.WithValue(r.Context(), "token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
