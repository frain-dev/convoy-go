package convoy_go

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Convoy struct {
	options Options
}

type Options struct {
	HTTPClient  HTTPClient
	APIEndpoint string
	APIUsername string
	APIPassword string
}

func New() *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIEndpoint: retrieveUrlFromEnv(),
		APIUsername: retrieveUsernameFromEnv(),
		APIPassword: retrievePasswordFromEnv(),
	}
	return &Convoy{
		options: options,
	}
}

func NewWithCredentials(url, username, password string) *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIEndpoint: url,
		APIUsername: username,
		APIPassword: password,
	}
	return &Convoy{
		options: options,
	}
}

func retrieveUrlFromEnv() string {
	url := os.Getenv("CONVOY_URL")
	if url == "" || len(url) == 0 {
		log.Println("Unable to retrieve Convoy URL")
	}
	return url
}

func retrieveUsernameFromEnv() string {
	username := os.Getenv("CONVOY_API_USERNAME")
	if username == "" || len(username) == 0 {
		log.Println("Unable to retrieve Convoy API username")
	}
	return username
}

func retrievePasswordFromEnv() string {
	password := os.Getenv("CONVOY_API_PASSWORD")
	if password == "" || len(password) == 0 {
		log.Println("Unable to retrieve Convoy API password")
	}
	return password
}
