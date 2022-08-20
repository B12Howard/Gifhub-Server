package services

import (
	"encoding/json"
	vidprocessing "gifconverter/services/vid-processing"
	"gifconverter/shared/utility/delete_file"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ExtractByUrlData struct {
	Video string
	Start string
	Dur   int // seconds
}

type ExtractByUrlDataStartEnd struct {
	Video    string
	Start    string
	End      string
	WsUserID string
	Id       int
}

const (
	GCPBucket = "created-gifs"
)

func ServeExtractByUrlGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ytURL := r.URL.Query().Get("video")
		start := r.URL.Query().Get("start")
		dur := r.URL.Query().Get("dur")

		durI, err := strconv.Atoi(dur)
		if err != nil {
			panic(err)
		}

		id := uuid.New()
		fileName := id.String()
		fullPath := vidprocessing.OutDir + fileName + ".gif"

		file, err := vidprocessing.ConvertToGifByUrl(ytURL, start, durI, fullPath)
		if err != nil {
			panic(err)
		}

		// TODO pass back the url of the file OR encode as base64
		// TODO after convert to blob delete file?
		render.Status(r, http.StatusFound)
		render.JSON(w, r, file)

		http.Redirect(w, r, "/gif/"+fileName, 302)
	}
}

// Synchronous version - keeps http connection open; not useful but leaving it here for example usage
func ServeExtractByUrlSynchronous() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data ExtractByUrlData

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			panic(err)
		}

		id := uuid.New()
		fileName := id.String()
		fullPath := vidprocessing.OutDir + fileName + ".gif"

		file, err := vidprocessing.ConvertToGifByUrl(data.Video, data.Start, data.Dur, fullPath)
		if err != nil {
			panic(err)
		}

		rmvError := delete_file.RemoveFileFromDirectory(fullPath)
		if rmvError != nil {
			panic(rmvError)
		}

		render.Status(r, http.StatusFound)
		render.JSON(w, r, file)

	}
}

// Asynchronous version
func ServeExtractByUrl(hub *Hub, db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data ExtractByUrlDataStartEnd

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			panic(err)
		}

		go convertToGifConcurrentStartEnd(data, hub, data.WsUserID, db)

		if err != nil {
			panic(err)
		}

		render.Status(r, http.StatusFound)
		render.JSON(w, r, "procressing file")

	}
}

// Only use when processing more than 1 file at a time
// Chunking a video took longer than processing the whole thing at once
func ServeExtractByUrlWithConcurrency() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data ExtractByUrlData
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

			go convertToGifConcurrent(i, c, data, data.Start, data.Dur)
		}

		wg.Wait()

		for i := 0; i < elementCount; i++ {
			s = append(s, <-c)
		}

		render.Status(r, http.StatusFound)
		render.JSON(w, r, s)
	}
}

func convertToGifConcurrentStartEnd(data ExtractByUrlDataStartEnd, hub *Hub, wsId string, db *sql.DB) {

	id := uuid.New()
	fileName := id.String()
	fullPath := vidprocessing.OutDir + fileName + ".gif"
	start := time.Now()
	_, errProcessing := vidprocessing.ConvertToGifByUrlByStartEnd(data.Video, data.Start, data.End, fullPath)

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

	EmitToSpecificClient(hub, socketEventResponse, wsId)

	return
}

func convertToGifConcurrent(i int, c chan []byte, data ExtractByUrlData, choppedStart string, choppedEnd int) {
	defer wg.Done()

	id := uuid.New()
	fileName := id.String()
	fullPath := vidprocessing.OutDir + fileName + ".gif"

	file, err := vidprocessing.ConvertToGifByUrl(data.Video, choppedStart, choppedEnd, fullPath)
	if err != nil {
		panic(err)
	}

	c <- file

	return
}
