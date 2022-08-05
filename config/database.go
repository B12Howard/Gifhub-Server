package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type MyHandler struct {
	db *sql.DB
}

const (
	host     = ""
	port     = ""
	user     = ""
	password = ""
	dbname   = ""
)

func NewDb() *sql.DB {
	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal((err))
	}

	return db
}
