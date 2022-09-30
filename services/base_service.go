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

type GCPCloudStorageConfig struct {
	GCPCLOUDSTORAGE struct {
		Type                        string `json:"type"`
		Project_id                  string `json:"project_id"`
		Private_key_id              string `json:"private_key_id"`
		Private_key                 string `json:"private_key"`
		Client_email                string `json:"client_email"`
		Client_id                   string `json:"client_id"`
		Auth_uri                    string `json:"auth_uri"`
		Token_uri                   string `json:"token_uri"`
		Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
		Client_x509_cert_url        string `json:"client_x509_cert_url"`
	} `json:GCPCLOUDSTORAGE`
}
