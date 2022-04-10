package auth

import (
	"context"
	"net/http"
	"path/filepath"

	"strings"

	firebase "firebase.google.com/go"
	"github.com/go-chi/render"
	"google.golang.org/api/option"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
		opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			return
		}

		auth, err := app.Auth(context.Background())
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			return
		}

		header := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
