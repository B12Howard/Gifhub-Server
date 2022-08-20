package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

func GetIndexHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print()
		var todos []string
		var id int
		var uid string
		var createdat sql.NullTime
		var usertypeid int
		var updatedat sql.NullTime
		var deletedat sql.NullTime
		rows, err := db.Query("SELECT * FROM users")

		defer rows.Close()

		if err != nil {
			log.Fatalln(err)
			render.JSON(w, r, ("An Error Occured"))
		}

		for rows.Next() {
			fmt.Print(rows.Scan(&id, &createdat, &usertypeid, &uid, &updatedat, &deletedat))
			rows.Scan(&createdat, &usertypeid, &id, &uid, &updatedat, &deletedat)
			todos = append(todos, uid, strconv.Itoa(id))
		}

		// return list
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, todos)
	}
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"stuff": "post"})
}
func PutHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"stuff": "put"})
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"stuff": "delete"})
}
