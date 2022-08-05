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

type UserUsageQuery struct {
	Uid       string       `json:"uid"`
	StartDate sql.NullTime `json:"startDate"`
	Timespan  int          `json:"timespan"` // days
	// Pagination Pagination   `json:"pagination"`
}

type UserUsage struct {
	Id        int          `json:"id"`
	Uid       string       `json:"uid"`
	Duration  int          `json:"duration"`
	Createdat sql.NullTime `json:"createdat"`
}

type UserUsageRes struct {
	Usage         []UserUsage `json:"usage"`
	TotalDuration int         `json:"totalduration"`
}

func GetUserUsage(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data UserUsageQuery
		var payload UserUsageRes
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		errDecode := decoder.Decode(&data)

		if errDecode != nil {
			log.Fatalln(errDecode)
			render.JSON(w, r, ("Bad request. Invalid data sent"))
		}
		now := time.Now().UTC()
		if data.Timespan == 0 {
			data.Timespan = defaultTimespan
		}

		timeLowerBound := now.AddDate(0, 0, -data.Timespan)

		rows, errRows := db.Query(`SELECT Usage.id, Usage.duration, Usage.createdat FROM usage Usage INNER JOIN users Users ON Users.id=Usage.uid WHERE Users.uid=$1 AND Usage.createdat between $2 AND $3`, data.Uid, timeLowerBound, now)
		total := db.QueryRow(`SELECT SUM(UsageQuery.d) AS total FROM (SELECT Usage.id, Usage.duration AS d, Usage.createdat FROM usage Usage INNER JOIN users Users ON Users.id=Usage.uid WHERE Users.uid=$1 AND Usage.createdat between $2 AND $3) UsageQuery`, data.Uid, timeLowerBound, now)

		if errRows != nil {
			log.Fatalln(errRows)
			render.JSON(w, r, ("Error fetching rows"))
		}

		for rows.Next() {
			userUsage := UserUsage{}
			rows.Scan(&userUsage.Id, &userUsage.Duration, &userUsage.Createdat)
			payload.Usage = append(payload.Usage, userUsage)
		}

		total.Scan(&payload.TotalDuration)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, payload)
	}
}
