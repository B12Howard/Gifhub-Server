package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type UserFilePagination struct {
	RowCount int              `json:"rowCount"`
	LastId   int              `json:"lastId"`
	LastDate sql.NullTime     `json:"lastDate"`
	Records  []UserFileRecord `json:"records"`
	Next     bool             `json:"next"`
}

type UserFileRecord struct {
	Id        int          `json:"id"`
	Url       string       `json:"url"`
	Createdat sql.NullTime `json:"createdat"`
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
			// log.Fatalln(errDecode)
			render.JSON(w, r, ("Missing number of pagination rows"))
		}

		if data.LastDate.Valid == false {
			rows, errRows = db.Query(`SELECT id, url, createdat FROM userfiles ORDER BY createdat DESC, id DESC FETCH FIRST $1 ROWS ONLY`, data.RowCount)

		} else if data.Next {
			rows, errRows = db.Query(`SELECT id, url, createdat FROM userfiles WHERE (createdat, id) < ($1, $2) ORDER BY createdat DESC, id DESC FETCH FIRST $3 ROWS ONLY`, data.LastDate.Time, data.LastId, data.RowCount)

		} else {
			rows, errRows = db.Query(`SELECT id, url, createdat FROM userfiles WHERE (createdat, id) > ($1, $2) ORDER BY createdat DESC, id ASC FETCH FIRST $3 ROWS ONLY`, data.LastDate.Time, data.LastId, data.RowCount)
		}

		if errRows != nil {
			// log.Fatalln(errRows)
			render.JSON(w, r, ("No records found"))
		}

		for rows.Next() {
			userFile := UserFileRecord{}
			rows.Scan(&userFile.Id, &userFile.Url, &userFile.Createdat)
			data.Records = append(data.Records, userFile)
		}
		fmt.Print(&data.Records[0])
		data.LastId = data.Records[len(data.Records)-1].Id
		data.LastDate = data.Records[len(data.Records)-1].Createdat

		render.Status(r, http.StatusOK)
		render.JSON(w, r, data)
	}
}
