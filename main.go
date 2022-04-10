package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"gifconverter/config"
	"gifconverter/router"
)

type MyHandler struct {
	db *sql.DB
}

type Todo struct {
	item string
	id   int
}

func main() {
	fmt.Println("Loading...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	db := config.NewDb()

	router := router.NewRoutes(db)

	fmt.Println("Ready!")

	http.ListenAndServe(":5020", router)
	log.Fatalln("connected to port 5000")

}
