package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// WebhookRequest represents the request payload received from the webhook
type WebhookRequest struct {
	Event     string      `json:"event"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature"`
}

var (
	subscribersMutex sync.Mutex
	subscribers      = make(map[string]chan<- WebhookRequest)
)

// Subscribe adds a new subscriber to the webhook
func Subscribe(w http.ResponseWriter, r *http.Request) {
	subscriberID := mux.Vars(r)["id"]

	// Create a channel for the subscriber to receive webhook events
	eventChannel := make(chan WebhookRequest)

	subscribersMutex.Lock()
	// Add the subscriber to the subscribers map
	subscribers[subscriberID] = eventChannel
	subscribersMutex.Unlock()

	// Send a success response
	w.WriteHeader(http.StatusOK)
}

// Unsubscribe removes a subscriber from the webhook
func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	subscriberID := mux.Vars(r)["id"]

	subscribersMutex.Lock()
	// Remove the subscriber from the subscribers map
	delete(subscribers, subscriberID)
	subscribersMutex.Unlock()

	// Send a success response
	w.WriteHeader(http.StatusOK)
}

// WebhookHandler handles incoming webhook requests
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a WebhookRequest object
	var webhookReq WebhookRequest
	err := json.NewDecoder(r.Body).Decode(&webhookReq)
	if err != nil {
		log.Println("Error decoding webhook request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subscribersMutex.Lock()
	defer subscribersMutex.Unlock()

	// Send the webhook request to all subscribers
	for _, subscriber := range subscribers {
		go func(subscriber chan<- WebhookRequest) {
			// Send the webhook request to the subscriber
			subscriber <- webhookReq
		}(subscriber)
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/webhook/subscribe/{id}", Subscribe).Methods("POST")
	r.HandleFunc("/webhook/unsubscribe/{id}", Unsubscribe).Methods("POST")
	r.HandleFunc("/webhook", WebhookHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
