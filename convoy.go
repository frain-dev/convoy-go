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
	HTTPClient HTTPClient

	APIUsername string

	APIPassword string
}

func New() *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIUsername: retrieveUsernameFromEnv(),
		APIPassword: retrievePasswordFromEnv(),
	}
	return &Convoy{
		options: options,
	}
}

func NewWithCredentials(username, password string) *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIUsername: username,
		APIPassword: password,
	}
	return &Convoy{
		options: options,
	}
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
