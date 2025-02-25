package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	// "fmt"

	// customMiddleware "kithli-api/router/custom_middleware"
	"kithli-api/firebase" // Import AuthMiddleware package
	"kithli-api/services"

	// "log"
	// "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	// "github.com/gorilla/websocket"
)

// var FirebaseAuth *auth.Client

// func setupFirebaseAuth() {
// 	app := FirebaseApp
// 	client, err := app.Auth(context.Background())
// 	if err != nil {
// 		log.Fatalf("Error initializing Firebase Auth client: %v", err)
// 	}
// 	FirebaseAuth = client
// }
type ContextKey string

const UserIDKey ContextKey = "user_id"

// AuthMiddleware verifies Firebase ID tokens and extracts the user ID
func AuthMiddleware(authClient *firebase.FirebaseClient) func(http.Handler) http.Handler {
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
			token, err := authClient.FirebaseAuth.VerifyIDToken(r.Context(), tokenString)
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

func NewRoutes(router *chi.Mux, db *sql.DB, firebaseClient *firebase.FirebaseClient) *chi.Mux {
	hub := services.NewHub()
	go hub.Run()

	router.Use(middleware.Logger)
	log.Println("Middleware Logger Enabled")

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://localhost:3000"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowedOrigins: []string{"https://app.kithli.com"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// router.Route("/ws/{userId}", func(router chi.Router) {
	// 	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 		var upgrader = websocket.Upgrader{
	// 			ReadBufferSize:  1024,
	// 			WriteBufferSize: 1024,
	// 		}
	// 		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// 		userId := chi.URLParam(r, "userId")
	// 		fmt.Println("userId", userId)
	// 		connection, err := upgrader.Upgrade(w, r, nil)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return
	// 		}

	// 		services.CreateNewSocketUser(hub, connection, userId)

	// 	})
	// 	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
	// 		userId := chi.URLParam(r, "userId")
	// 		var socketEventResponse services.SocketEventStruct
	// 		socketEventResponse.EventName = "message response"
	// 		socketEventResponse.EventPayload = map[string]interface{}{
	// 			"username": "usernamestuff",
	// 			"message":  "file is complete",
	// 			"userID":   userId,
	// 		}
	// 		services.EmitToSpecificClient(hub, socketEventResponse, userId)

	// 	})
	// })

	
	router.Group(func(router chi.Router) {
		router.Route("/create-user", func(router chi.Router) {
			router.Post("/", services.CreateUserHandler(db, firebaseClient))
		})
		router.Route("/create-member", func(router chi.Router) {
			router.Post("/", services.CreateMemberHandler(db))
		})
		// router.Use(Auth(firebaseClient))
		router.Route("/getUser", func(router chi.Router) {
			router.Post("/", services.GetUser(db))
			// Example if not passing in to Auth
			// router.With(AuthMiddleware(firebaseClient)).Post("/create-user", CreateUserHandler(firebaseClient))
			// router.Post("/getGifs", services.GetUserGifs(db))
			// router.Delete("/deleteGif", services.DeleteGifById(db))
			// router.Post("/getUsage", services.GetUserUsage(db))
		})

		// router.Route("/useConverter", func(router chi.Router) {
		// 	router.Post("/convertVIdeosToGifsStitchTogether", services.ConvertVIdeosToGifsStitchTogether())
		// 	router.Post("/convertVideoToGif", services.ConvertVideoToGif(hub, db))
		// })

		// router.Route("/getSignedUrlGif", func(router chi.Router) {
		// 	router.Post("/", services.GetUserImage(db))
		// })

		router.Route("/check-membership", func(router chi.Router) {
			router.Post("/", services.CheckUserMembershipHandler(db))
		})

		router.Route("/get-user-data", func(router chi.Router) {
			router.With(AuthMiddleware(firebaseClient)).Get("/{external_id}", services.GetUserDataHandler(db, firebaseClient))
		})

	})
	return router
}
