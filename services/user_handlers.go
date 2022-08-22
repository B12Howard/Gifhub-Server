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
	Name       string
	Filelimit  int
	Usagelimit int
	Id         int
	Createdat  sql.NullTime
	Updatedat  sql.NullTime
}

type UserRes struct {
	id         int
	uid        string
	usertypeid int
	createdat  sql.NullTime
	ut         UserTypes
}

func GetUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data UserQuery
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		errDecode := decoder.Decode(&data)

		if errDecode != nil {
			log.Fatalln(errDecode)
			render.JSON(w, r, ("Missing uid"))
		}

		retrievedUser := &UserRes{}
		retrievedUserType := &UserTypes{}

		row := db.QueryRow(`SELECT Users.id, Users.createdat, UserTypes.name,  UserTypes.filelimit,  UserTypes.createdat, uid FROM users Users INNER JOIN usertypes UserTypes ON Users.usertypeid=UserTypes.id WHERE uid = $1`, data.Uid)
		err := row.Scan(&retrievedUser.id, &retrievedUser.createdat, &retrievedUserType.Name, &retrievedUserType.Filelimit, &retrievedUserType.Createdat, &retrievedUser.uid)

		if err != nil {
			log.Fatalln(err)
			render.JSON(w, r, ("No User found with uid " + data.Uid))
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