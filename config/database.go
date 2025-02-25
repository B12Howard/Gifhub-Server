package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type MyHandler struct {
	db *sql.DB
}

type DBConfig struct {
	DB struct {
		Host     string `json:"HOST"`
		Port     int    `json:"PORT"`
		User     string `json:"USER"`
		Password string `json:"PASSWORD"`
		Dbname   string `json:"DBNAME"`
	} `json:DB`
}

func NewDb() *sql.DB {
	var dbConfig DBConfig

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&dbConfig)

	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.DB.Host, dbConfig.DB.Port, dbConfig.DB.User, dbConfig.DB.Password, dbConfig.DB.Dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err)
	}

	return db
}
