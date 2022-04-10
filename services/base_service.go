package services

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

func GetIndexHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var todos []string
		var id int
		var item string
		rows, err := db.Query("SELECT * FROM todos")

		defer rows.Close()

		if err != nil {
			log.Fatalln(err)
			render.JSON(w, r, ("An Error Occured"))
		}

		for rows.Next() {
			rows.Scan(&item, &id)
			todos = append(todos, item)
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
