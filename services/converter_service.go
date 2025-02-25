package services

import (
	"database/sql"
	"encoding/json"
	vidprocessing "kithli-api/services/vid-processing"
	"kithli-api/shared/utility/delete_file"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type VideoToGifByDuration struct {
	Video string
	Start string
	Dur   int // seconds
}

type VideoToGifByStartEnd struct {
	Video    string
	Start    string
	End      string
	WsUserID string
	Id       int
}

// ConvertVideoToGif takes in a VideoToGifByDuration and closes the http connection while converting the video to gif.
// Calls background task is completeConvertToGifByStartEnd()
// Returns status 200 and a message to the user letting them know we have received their request.
func ConvertVideoToGif(hub *Hub, db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data VideoToGifByStartEnd
		retrievedUser := &UserRoleLimits{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			log.Println(err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Error: Bad Video Parameters")
			return
		}

		row := db.QueryRow(`SELECT u.id, ut.max_gif_time, ut.file_size_limit, ut.usage_limit FROM users u JOIN user_types ut on ut.id=u.user_type_id WHERE uid = $1 `, data.WsUserID)
		rowErr := row.Scan(&retrievedUser.id, &retrievedUser.max_gif_time, &retrievedUser.file_size_limit, &retrievedUser.usage_limit)

		if rowErr != nil {
			log.Println(rowErr)
		}

		start := strings.Split(data.Start, ":")
		end := strings.Split(data.End, ":")

		startHour, _ := strconv.Atoi(start[0])
		startMin, _ := strconv.Atoi(start[1])
		startSec, _ := strconv.Atoi(start[2])
		endHour, _ := strconv.Atoi(end[0])
		endMin, _ := strconv.Atoi(end[1])
		endSec, _ := strconv.Atoi(end[2])

		t1 := time.Date(1984, time.November, 3, startHour, startMin, startSec, 0, time.UTC)
		t2 := time.Date(1984, time.November, 3, endHour, endMin, endSec, 0, time.UTC)

		if t2.Sub(t1).Seconds() <= float64(0) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Start time has to be before the end time")
			return
		}

		if t2.Sub(t1).Seconds() > float64(retrievedUser.max_gif_time) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Gif duration is longer than your subscription limit")
			return
		}

		go completeConvertToGifByStartEnd(data, hub, data.WsUserID, db, retrievedUser.id)

		if err != nil {
			log.Println(err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Error While Converting")
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "Procressing File")

	}
}

// ConvertVIdeosToGifsStitchTogether takes in an array of VideoToGifByDuration
// Calls completeConvertVideosToGifs then concats them together.
// This may not work as expected HL
func ConvertVIdeosToGifsStitchTogether() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data []VideoToGifByDuration
		err := decoder.Decode(&data)

		if err != nil {
			log.Println(err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Error: Bad Video Parameters")
			return
		}

		elementCount := 2

		if elementCount == 0 {
			elementCount = 1
		}

		s := make([][]byte, 0, elementCount)
		c := make(chan []byte, elementCount)
		// completed := 0

		for i := 0; i < elementCount; i++ {
			wg.Add(1)

			go completeConvertVideosToGifs(i, c, data[i], data[i].Start, data[i].Dur)
		}

		wg.Wait()

		for i := 0; i < elementCount; i++ {
			s = append(s, <-c)
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, s)
	}
}

// completeConvertToGifByStartEnd takes a VideoToGifByStartEnd and creates a gif via ConvertToGifCutByStartEnd
// If successful the gif is saved to GCP.
// The duration/ usage time is saved to Postgres
// If successful send the user a message via websocket
func completeConvertToGifByStartEnd(data VideoToGifByStartEnd, hub *Hub, wsId string, db *sql.DB, retrievedUserId int) {

	id := uuid.New()
	fileName := id.String()
	fullPath := vidprocessing.OutDir + fileName + ".gif"
	start := time.Now()
	_, errProcessing := vidprocessing.ConvertToGifCutByStartEnd(data.Video, data.Start, data.End, fullPath)

	var socketEventResponse SocketEventStruct
	socketEventResponse.EventName = "message response"

	if errProcessing != nil {
		log.Panic(errProcessing)
		socketEventResponse.EventPayload = map[string]interface{}{
			"username": fileName,
			"message":  "Error processing gifs",
			"userID":   data.WsUserID,
		}

		EmitToSpecificClient(hub, socketEventResponse, wsId)

		return
	}

	objectUrl := "https://storage.cloud.google.com/" + GCPBucket + "/" + fileName

	f, _ := os.Open(fullPath)
	defer f.Close()

	_, err := db.Exec("INSERT INTO user_files (url, created_at, uid) VALUES ($1, $2, $3)", objectUrl, time.Now().UTC(), retrievedUserId)

	if err != nil {
		log.Println(err)
		socketEventResponse.EventPayload = map[string]interface{}{
			"username": fileName,
			"message":  "Error saving file url",
			"userID":   data.WsUserID,
		}

		EmitToSpecificClient(hub, socketEventResponse, wsId)

		return
	}

	errFileUpload := FileUpload(GCPBucket, f, fileName)

	if errFileUpload != nil {
		log.Println(errFileUpload)

		socketEventResponse.EventPayload = map[string]interface{}{
			"username": fileName,
			"message":  "Error uploading file",
			"userID":   data.WsUserID,
		}

		EmitToSpecificClient(hub, socketEventResponse, wsId)

		return
	}

	_, errSaveFileUrl := db.Exec("INSERT INTO usage (uid, duration, created_at) VALUES ($1, $2, $3)", data.Id, math.Round(time.Now().Sub(start).Seconds()), time.Now().UTC())

	if errFileUpload != nil {
		log.Println(errSaveFileUrl)

		socketEventResponse.EventPayload = map[string]interface{}{
			"username": fileName,
			"message":  "Error saving uploaded file",
			"userID":   data.WsUserID,
		}

		EmitToSpecificClient(hub, socketEventResponse, wsId)

		return
	}

	rmvError := delete_file.RemoveFileFromDirectory(fullPath)
	if rmvError != nil {
		log.Println(rmvError)
	}

	EmitToSpecificClient(hub, socketEventResponse, wsId)

	return
}

// completeConvertVideosToGifs takes in a VideoToGifByDuration and pushes the completed gif to a channel
func completeConvertVideosToGifs(i int, c chan []byte, data VideoToGifByDuration, choppedStart string, choppedEnd int) {
	defer wg.Done()

	id := uuid.New()
	fileName := id.String()
	fullPath := vidprocessing.OutDir + fileName + ".gif"

	file, err := vidprocessing.ConvertToGifCutByDuration(data.Video, choppedStart, choppedEnd, fullPath)
	if err != nil {
		log.Panic(err)
	} else {
		c <- file
	}

	rmvError := delete_file.RemoveFileFromDirectory(fullPath)

	if rmvError != nil {
		log.Println(rmvError)
	}

	return
}
