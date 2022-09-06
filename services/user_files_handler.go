package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type UserFilePagination struct {
	RowCount int          `json:"rowCount"`
	LastId   int          `json:"lastId"`
	LastDate sql.NullTime `json:"lastDate"`
	Records  []UserFile   `json:"records"`
	Next     bool         `json:"next"`
}

type UserFile struct {
	Id         int          `json:"id"`
	Url        string       `json:"url"`
	Created_at sql.NullTime `json:"created_at"`
}

// GetUserGifs takes in a UserFilePagination and queries the `userfiles` table for gifs converted and saved to the remote storage (GCP Cloud Storage).
// Uses keyset pagination.
// Returns a UserFilePagination
// TODO Should only get the past 24 hours because 24 hours is the time limit the converted files live in the remote storage
func GetUserGifs(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data UserFilePagination
		var rows *sql.Rows
		var errRows error

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		errDecode := decoder.Decode(&data)

		if errDecode != nil {
			log.Println(errDecode)
			render.JSON(w, r, ("Missing number of pagination rows"))

			return
		}

		if data.LastDate.Valid == false {
			rows, errRows = db.Query(`SELECT id, url, created_at FROM user_files ORDER BY created_at DESC, id DESC FETCH FIRST $1 ROWS ONLY`, data.RowCount)

		} else if data.Next {
			rows, errRows = db.Query(`SELECT id, url, created_at FROM user_files WHERE (created_at, id) < ($1, $2) ORDER BY created_at DESC, id DESC FETCH FIRST $3 ROWS ONLY`, data.LastDate.Time, data.LastId, data.RowCount)

		} else {
			rows, errRows = db.Query(`SELECT id, url, created_at FROM user_files WHERE (created_at, id) > ($1, $2) ORDER BY created_at DESC, id ASC FETCH FIRST $3 ROWS ONLY`, data.LastDate.Time, data.LastId, data.RowCount)
		}

		if errRows != nil {
			log.Println(errRows)
			render.JSON(w, r, ("No records found"))

			return
		}

		for rows.Next() {
			userFile := UserFile{}
			rows.Scan(&userFile.Id, &userFile.Url, &userFile.Created_at)
			data.Records = append(data.Records, userFile)
		}
		fmt.Print(&data.Records[0])
		data.LastId = data.Records[len(data.Records)-1].Id
		data.LastDate = data.Records[len(data.Records)-1].Created_at

		render.Status(r, http.StatusOK)
		render.JSON(w, r, data)
	}
}
