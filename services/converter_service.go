package services

import (
	"encoding/json"
	vidprocessing "gifconverter/services/vid-processing"
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
}

func ServeExtractByUrlGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ytURL := r.URL.Query().Get("video")
		start := r.URL.Query().Get("start")
		dur := r.URL.Query().Get("dur")

		// startI, err := strconv.Atoi(start)
		// if err != nil {
		// 	panic(err)
		// }

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

// Synchronous version
// func ServeExtractByUrl() func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var data ExtractByUrlData

// 		decoder := json.NewDecoder(r.Body)
// 		err := decoder.Decode(&data)

// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("Received args: %v %v %v\n", data.Video, data.Start, data.Dur)

// 		id := uuid.New()
// 		fileName := id.String()
// 		fullPath := vidprocessing.OutDir + fileName + ".gif"

// 		file, err := vidprocessing.ConvertToGifByUrl(data.Video, data.Start, data.Dur, fullPath)
// 		if err != nil {
// 			panic(err)
// 		}

// 		rmvError := delete_file.RemoveFileFromDirectory(fullPath)
// 		if rmvError != nil {
// 			panic(rmvError)
// 		}

// 		render.Status(r, http.StatusFound)
// 		render.JSON(w, r, file)

// 	}
// }

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
			// choppedStart := data.Start + i*parts
			// if i == 0 {
			// 	choppedStart = data.Start
			// }

			// choppedEnd := choppedStart + parts
			// if data.Dur < choppedEnd {
			// 	choppedEnd = data.Dur - completed
			// }
			// completed = completed + parts
			go convertToGifConcurrent(i, c, data, data.Start, data.Dur)
		}

		wg.Wait()

		for i := 0; i < elementCount; i++ {
			s = append(s, <-c)
		}

		//http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/
		// outGif := &gif.GIF{}

		// for _, name := range s {
		// 	f, _ := os.Open(vidprocessing.OutDir + name + ".gif")
		// 	inGif, _ := gif.Decode(f)
		// 	f.Close()

		// 	outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
		// 	outGif.Delay = append(outGif.Delay, 0)
		// }

		// id := uuid.New().String()

		// // save to out.gif
		// f, _ := os.OpenFile(vidprocessing.OutDir+"final_"+id+".gif", os.O_WRONLY|os.O_CREATE, 0600)
		// defer f.Close()
		// gif.EncodeAll(f, outGif)

		render.Status(r, http.StatusFound)
		render.JSON(w, r, s)

		// http.Redirect(w, r, "/gif/"+fileName, 302)
	}
}

func convertToGifConcurrentStartEnd(data ExtractByUrlDataStartEnd, hub *Hub, userID string, db *sql.DB) {

	id := uuid.New()
	fileName := id.String()
	fullPath := vidprocessing.OutDir + fileName + ".gif"
	start := time.Now()
	_, err := vidprocessing.ConvertToGifByUrlByStartEnd(data.Video, data.Start, data.End, fullPath)

	if err != nil {
		panic(err)
	}
	f, _ := os.Open(fullPath)
	defer f.Close()

	_, err = db.Exec("INSERT INTO usage (uid, duration, createdat) VALUES ($1, $2, $3)", userID, math.Round(time.Now().Sub(start).Seconds()), time.Now().UTC())

	if err != nil {
		panic(err)
	}

	FileUpload("created-gifs", f, fileName)

	var socketEventResponse SocketEventStruct
	socketEventResponse.EventName = "message response"
	socketEventResponse.EventPayload = map[string]interface{}{
		"username": "usernamestuff",
		"message":  "file is complete",
		"userID":   userID,
	}

	EmitToSpecificClient(hub, socketEventResponse, userID)

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
