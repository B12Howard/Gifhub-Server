package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type MyHandler struct {
	db *sql.DB
}

func NewDb() *sql.DB {
	connectionString := os.Getenv("CONNECTIONSTRING")
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal((err))
	}

	return db
}
