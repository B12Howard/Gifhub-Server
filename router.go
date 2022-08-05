package main

import (
	"database/sql"
	"fmt"
	"gifconverter/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

func NewRoutes(router *chi.Mux, db *sql.DB) *chi.Mux {
	hub := services.NewHub()
	go hub.Run()

	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Route("/getUser", func(router chi.Router) {
		router.Post("/", services.GetUser(db))
		router.Post("/getGifs", services.GetUserGifs(db))
		router.Post("/getUsage", services.GetUserUsage(db))
	})
	router.Route("/useConverter", func(router chi.Router) {
		router.Post("/", services.ServeExtractByUrl(hub, db))
		router.Post("/withConcurrency", services.ServeExtractByUrlWithConcurrency())

	})
	router.Route("/concurrency", func(router chi.Router) {

		router.Get("/", services.ServeConcurrency())

	})
	router.Route("/ws/{userId}", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var upgrader = websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
			userId := chi.URLParam(r, "userId")
			fmt.Println("userId", userId)
			connection, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			services.CreateNewSocketUser(hub, connection, userId)

		})
		router.Post("/", func(w http.ResponseWriter, r *http.Request) {
			userId := chi.URLParam(r, "userId")
			var socketEventResponse services.SocketEventStruct
			socketEventResponse.EventName = "message response"
			socketEventResponse.EventPayload = map[string]interface{}{
				"username": "usernamestuff",
				"message":  "file is complete",
				"userID":   userId,
			}
			services.EmitToSpecificClient(hub, socketEventResponse, userId)

		})
	})
	router.Route("/", func(router chi.Router) {
		router.Get("/", services.GetIndexHandler(db))
		router.Post("/", services.PostHandler)
		router.Put("/", services.PutHandler)
		router.Delete("/", services.DeleteHandler)
	})
	return router
}
