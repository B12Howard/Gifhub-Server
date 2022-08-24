package services

import (
	"encoding/json"
	vidprocessing "gifconverter/services/vid-processing"
	"gifconverter/shared/utility/delete_file"
	"math"
	"net/http"
	"os"
	"time"

	"database/sql"

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

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			panic(err)
		}

		go completeConvertToGifByStartEnd(data, hub, data.WsUserID, db)

		if err != nil {
			panic(err)
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "procressing file")

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
			panic(err)
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
func completeConvertToGifByStartEnd(data VideoToGifByStartEnd, hub *Hub, wsId string, db *sql.DB) {

	id := uuid.New()
	fileName := id.String()
	fullPath := vidprocessing.OutDir + fileName + ".gif"
	start := time.Now()
	_, errProcessing := vidprocessing.ConvertToGifCutByStartEnd(data.Video, data.Start, data.End, fullPath)

	if errProcessing != nil {
		panic(errProcessing)
	}

	f, _ := os.Open(fullPath)
	defer f.Close()
	objectUrl := "https://storage.cloud.google.com/" + GCPBucket + fileName
	_, err := db.Exec("INSERT INTO userfiles (url, createdat, uid) VALUES ($1, $2, $3)", objectUrl, time.Now().UTC(), data.Id)

	if err != nil {
		panic(err)
	}

	errFileUpload := FileUpload(GCPBucket, f, fileName)

	if errFileUpload != nil {
		panic(errFileUpload)
	}

	_, errSaveFileUrl := db.Exec("INSERT INTO usage (uid, duration, createdat) VALUES ($1, $2, $3)", data.Id, math.Round(time.Now().Sub(start).Seconds()), time.Now().UTC())

	if errFileUpload != nil {
		panic(errSaveFileUrl)
	}

	var socketEventResponse SocketEventStruct
	socketEventResponse.EventName = "message response"
	socketEventResponse.EventPayload = map[string]interface{}{
		"username": "usernamestuff",
		"message":  "file is complete",
		"userID":   data.Id,
	}

	rmvError := delete_file.RemoveFileFromDirectory(fullPath)
	if rmvError != nil {
		panic(rmvError)
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
		panic(err)
	}

	c <- file

	rmvError := delete_file.RemoveFileFromDirectory(fullPath)
	if rmvError != nil {
		panic(rmvError)
	}

	return
}
