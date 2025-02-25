package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"kithli-api/config"
	"kithli-api/firebase"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type MyHandler struct {
	db *sql.DB
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Loading...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	db := config.NewDb()

	// fmt.Print("Hello world")

	// wd, a := os.Getwd()
	// if a != nil {
	// 	log.Fatalf("Failed to get working directory: %v", a)
	// }
	// fmt.Println("Current working directory:", wd)
	// _, r := os.Stat("./config/firebase.json")
	// if os.IsNotExist(r) {
	// 	log.Fatalf(".env file not found at ../config/firebase.json")
	// } else if r != nil {
	// 	log.Fatalf("Error accessing .env file: %v", r)
	// } else {
	// 	fmt.Println(".env file exists and is accessible.")
	// }

	firebaseClient, err := firebase.InitFirebase("./config/firebase.json")
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	log.Println("Running on port %s", port)

	router := chi.NewRouter()
	fmt.Println(firebaseClient)
	NewRoutes(router, db, firebaseClient)
	fmt.Println("Ready!")
	listenAndServePort := ":" + (port)
	http.ListenAndServe(listenAndServePort, router)
	fmt.Println("connected to port %s", port)

}
