package services

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/render"
)

var wg sync.WaitGroup

type concurrencyResult struct {
	mu   sync.Mutex
	data []int
}

func ServeConcurrency() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		elementCount := 10
		s := make([]string, 0, elementCount)
		c := make(chan int, elementCount)

		done := make(chan bool)

		for i := 0; i < elementCount; i++ {
			wg.Add(1)
			go concurrency(i, c, done)
		}

		wg.Wait()

		for i := 0; i < elementCount; i++ {
			s = append(s, strconv.Itoa(<-c)+"this!")
		}

		close(c)

		render.Status(r, http.StatusFound)
		render.JSON(w, r, s)
	}
}

func concurrency(i int, c chan int, done chan bool) {
	defer wg.Done()

	c <- i

	amt := time.Duration(rand.Intn(250))
	time.Sleep(time.Millisecond * amt)
	return
}
