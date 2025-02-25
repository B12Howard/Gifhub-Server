package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type UserQuery struct {
	Uid string `json:"uid"`
}

type UserTypes struct {
	Name            string
	File_size_limit int
	Usage_limit     int
	Id              int
	Created_at      sql.NullTime `json:"created_at"`
	Updated_at      sql.NullTime `json:"updated_at"`
}

type UserRes struct {
	id           int
	uid          string
	user_type_id int
	created_at   sql.NullTime
	ut           UserTypes
}

type UserRoleLimits struct {
	id              int
	max_gif_time    int
	file_size_limit int
	usage_limit     int
}

func GetUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data UserQuery
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		errDecode := decoder.Decode(&data)

		if errDecode != nil {
			log.Println(errDecode)
			render.JSON(w, r, ("Missing uid"))
		}

		retrievedUser := &UserRes{}
		retrievedUserType := &UserTypes{}

		row := db.QueryRow(`SELECT Users.id, Users.created_at, UserTypes.name,  UserTypes.file_size_limit,  UserTypes.created_at, uid FROM users Users INNER JOIN user_types UserTypes ON Users.user_type_id=UserTypes.id WHERE uid = $1`, data.Uid)
		err := row.Scan(&retrievedUser.id, &retrievedUser.created_at, &retrievedUserType.Name, &retrievedUserType.File_size_limit, &retrievedUserType.Created_at, &retrievedUser.uid)

		if err != nil {
			log.Println(err)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ("No User found with uid " + data.Uid))

			return
		}

		payload := map[string]interface{}{
			"id":  &retrievedUser.id,
			"uid": &retrievedUser.uid,
			"ut":  &retrievedUserType,
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, payload)
	}
}

func SetUserUsage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

// func GetUserUsageById()

// func SetUserUsage(w http.ResponseWriter, r *http.Request) {
// 	render.Status(r, http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"stuff": "post"})
// }
// func PutHandler(w http.ResponseWriter, r *http.Request) {
// 	render.Status(r, http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"stuff": "put"})
// }

// func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	render.Status(r, http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"stuff": "delete"})
// }
