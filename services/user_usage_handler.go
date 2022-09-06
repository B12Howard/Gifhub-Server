package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

const (
	defaultTimespan = 30
)

type UserUsagePagination struct {
	Uid       string       `json:"uid"`
	StartDate sql.NullTime `json:"startDate"`
	Timespan  int          `json:"timespan"` // days
	// Pagination Pagination   `json:"pagination"`
}

type UserUsage struct {
	Id         int          `json:"id"`
	Uid        string       `json:"uid"`
	Duration   int          `json:"duration"`
	Created_at sql.NullTime `json:"created_at"`
}

type UserUsageRes struct {
	Usage         []UserUsage `json:"usage"`
	TotalDuration int         `json:"totalduration"`
}

// TODO Implement User Usage CRUD

// GetUserUsage takes in a UserUsagePagination and queries the `usage` table.
// Returns a UserUsageRes
// TODO implement keyset pagination like GetUserGifs
func GetUserUsage(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data UserUsagePagination
		var payload UserUsageRes
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		errDecode := decoder.Decode(&data)

		if errDecode != nil {
			log.Println(errDecode)
			render.JSON(w, r, ("Bad request. Invalid data sent"))

			return
		}
		now := time.Now().UTC()
		if data.Timespan == 0 {
			data.Timespan = defaultTimespan
		}

		timeLowerBound := now.AddDate(0, 0, -data.Timespan)

		rows, errRows := db.Query(`SELECT Usage.id, Usage.duration, Usage.created_at FROM usage Usage INNER JOIN users Users ON Users.id=Usage.uid WHERE Users.uid=$1 AND Usage.created_at between $2 AND $3`, data.Uid, timeLowerBound, now)
		total := db.QueryRow(`SELECT SUM(UsageQuery.d) AS total FROM (SELECT Usage.id, Usage.duration AS d, Usage.created_at FROM usage Usage INNER JOIN users Users ON Users.id=Usage.uid WHERE Users.uid=$1 AND Usage.created_at between $2 AND $3) UsageQuery`, data.Uid, timeLowerBound, now)

		if errRows != nil {
			log.Println(errRows)
			render.JSON(w, r, ("Error fetching rows"))

			return
		}

		for rows.Next() {
			userUsage := UserUsage{}
			rows.Scan(&userUsage.Id, &userUsage.Duration, &userUsage.Created_at)
			payload.Usage = append(payload.Usage, userUsage)
		}

		total.Scan(&payload.TotalDuration)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, payload)
	}
}
