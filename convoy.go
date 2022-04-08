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
	GroupID     string
	APIKey      string
	APIEndpoint string
	APIUsername string
	APIPassword string
}

func New() *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		GroupID:     retrieveGroupIDFromEnv(),
		APIKey:      retrieveAPIKeyFromEnv(),
		APIEndpoint: retrieveURLFromEnv(),
		APIUsername: retrieveUsernameFromEnv(),
		APIPassword: retrievePasswordFromEnv(),
	}
	return &Convoy{
		options: options,
	}
}

func NewWithCredentials(url, groupID, username, password string) *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		GroupID:     groupID,
		APIEndpoint: url,
		APIUsername: username,
		APIPassword: password,
	}

	return &Convoy{
		options: options,
	}
}

func NewWithAPIKey(url, groupID, apiKey string) *Convoy {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		GroupID: groupID,
		APIKey:  apiKey,
	}

	return &Convoy{
		options: options,
	}
}

func retrieveGroupIDFromEnv() string {
	groupID := os.Getenv("CONVOY_GROUP_ID")
	if !isStringEmpty(groupID) {
		log.Println("Unable to retrieve Convoy groupID")
	}
	return groupID
}

func retrieveURLFromEnv() string {
	url := os.Getenv("CONVOY_URL")
	if isStringEmpty(url) {
		log.Println("Unable to retrieve Convoy URL")
	}
	return url
}

func retrieveUsernameFromEnv() string {
	username := os.Getenv("CONVOY_API_USERNAME")
	if isStringEmpty(username) {
		log.Println("Unable to retrieve Convoy API username")
	}
	return username
}

func retrievePasswordFromEnv() string {
	password := os.Getenv("CONVOY_API_PASSWORD")
	if isStringEmpty(password) {
		log.Println("Unable to retrieve Convoy API password")
	}
	return password
}

func retrieveAPIKeyFromEnv() string {
	apiKey := os.Getenv("CONVOY_API_KEY")
	if isStringEmpty(apiKey) {
		log.Println("Unable to retrieve api key")
	}
	return apiKey
}
