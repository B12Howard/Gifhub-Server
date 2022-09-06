package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"gifconverter/config"

	"github.com/go-chi/chi/v5"
)

type MyHandler struct {
	db *sql.DB
}

func main() {
	fmt.Println("Loading...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	db := config.NewDb()

	router := chi.NewRouter()
	NewRoutes(router, db)
	fmt.Println("Ready!")

	http.ListenAndServe(":5020", router)
	log.Fatalln("connected to port 5000")

}
