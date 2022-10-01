package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type FirebaseServiceAccountConfig struct {
	FIREBASE struct {
		Type                        string `json:"type"`
		Project_id                  string `json:"project_id"`
		Private_key_id              string `json:"private_key_id"`
		Private_key                 string `json:"private_key"`
		Client_email                string `json:"client_email"`
		Client_id                   string `json:"client_id"`
		Auth_uri                    string `json:"auth_uri"`
		Token_uri                   string `json:"token_uri"`
		Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
		Client_x509_cert_url        string `json:"client_x509_cert_url"`
	} `json:FIREBASE`
}

// Auth checks the request's jwt token for validity.
// Returns for invalid jwt tokens status 404
// Continues
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println("ex", exPath)
		viper.SetConfigName("config")
		viper.AddConfigPath(exPath + "/config")
		viperErr := viper.ReadInConfig()

		if viperErr != nil {
			fmt.Printf("Error reading config file, %s", viperErr)
		}

		var gcpConfig FirebaseServiceAccountConfig
		viperErr = viper.Unmarshal(&gcpConfig)
		b, _ := json.Marshal(gcpConfig.FIREBASE)
		opt := option.WithCredentialsJSON(b)
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
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
