package services

import (
	"sync"
	"time"
)

const (
	GCPBucket = "created-gifs"
)

var wg sync.WaitGroup

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)
