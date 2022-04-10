package router

import (
	"database/sql"
	"gifconverter/services"

	auth "gifconverter/router/custom_middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRoutes(db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
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
	router.Use(auth.Auth)

	router.Route("/converter", func(router chi.Router) {

		// router.Get("/", services.ServeExtractByUrl())
		router.Post("/", services.ServeExtractByUrl())
		router.Post("/concurrency", services.ServeExtractByUrlWithConcurrency())
		// router.Put("/converter", services.PutHandler)
		// router.Delete("/converter", services.DeleteHandler)
	})
	router.Route("/concurrency", func(router chi.Router) {

		router.Get("/", services.ServeConcurrency())
		// router.Post("/converter", services.PostHandler)
		// router.Put("/converter", services.PutHandler)
		// router.Delete("/converter", services.DeleteHandler)
	})
	router.Route("/", func(router chi.Router) {

		// router.Get("/asdf", services.GetIndexHandler(db))
		router.Get("/", services.GetIndexHandler(db))
		router.Post("/", services.PostHandler)
		router.Put("/", services.PutHandler)
		router.Delete("/", services.DeleteHandler)
	})
	return router
}
