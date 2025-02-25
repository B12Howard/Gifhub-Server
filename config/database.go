package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // For loading .env file
	_ "github.com/lib/pq"
)

type MyHandler struct {
	db *sql.DB
}

// type DBConfig struct {
// 	DB struct {
// 		Host     string `json:"HOST"`
// 		Port     int    `json:"PORT"`
// 		User     string `json:"USER"`
// 		Password string `json:"PASSWORD"`
// 		Dbname   string `json:"DBNAME"`
// 	} `json:DB`
// }

func NewDb() *sql.DB {
	// var dbConfig DBConfig
	// wd, wdErr := os.Getwd()

	// if wdErr != nil {
	// 	log.Fatalf("Failed to get working directory: %v", wdErr)
	// }
	// fmt.Print("Hello world")
	// wd, a := os.Getwd()
	// if a != nil {
	// 	log.Fatalf("Failed to get working directory: %v", a)
	// }
	// fmt.Println("Current working directory:", wd)
	// _, r := os.Stat(".env")
	// if os.IsNotExist(r) {
	// 	log.Fatalf(".env file not found at ../.env")
	// } else if r != nil {
	// 	log.Fatalf("Error accessing .env file: %v", r)
	// } else {
	// 	fmt.Println(".env file exists and is accessible.")
	// }


	// envPath := filepath.Join(wd, ".env")
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	fmt.Printf("Database connection string: host=%s port=%s user=%s dbname=%s\n", host, port, user, dbname)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Database connection established successfully!")
	return db
}
