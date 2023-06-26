package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// FileRequest represents a file request entity
type FileRequest struct {
	ID     string `json:"id"`
	FileID string `json:"fileId"`
}

var (
	requestsMutex sync.Mutex
	requests      = make(map[string]FileRequest)
)

// RequestFile creates a new file request
func RequestFile(w http.ResponseWriter, r *http.Request) {
	var fileRequest FileRequest
	// Parse the request body into the FileRequest object
	if err := json.NewDecoder(r.Body).Decode(&fileRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the file request
	fileRequest.ID = generateUniqueID()

	requestsMutex.Lock()
	// Save the file request
	requests[fileRequest.ID] = fileRequest
	requestsMutex.Unlock()

	// Return the created file request
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileRequest)
}

// SendFile sends the file to the hoster
func SendFile(w http.ResponseWriter, r *http.Request) {
	fileRequestID := mux.Vars(r)["id"]

	requestsMutex.Lock()
	// Check if the file request exists
	if fileRequest, ok := requests[fileRequestID]; ok {
		// Perform the file sending logic
		sendFileToHoster(fileRequest)

		// Remove the file request after sending the file
		delete(requests, fileRequestID)

		w.WriteHeader(http.StatusOK)
	} else {
		// File request not found
		w.WriteHeader(http.StatusNotFound)
	}
	requestsMutex.Unlock()
}

// sendFileToHoster is a mock function to simulate sending the file to the hoster
func sendFileToHoster(fileRequest FileRequest) {
	fmt.Println("Sending file with ID:", fileRequest.FileID, "to the hoster...")
	// Implement your logic to send the file to the hoster here
}

// generateUniqueID generates a unique ID for the file request
// You can replace this with your own unique ID generation logic
func generateUniqueID() string {
	// Implement your unique ID generation logic here
	return "unique-id"
}

func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/file-requests", RequestFile).Methods("POST")
	r.HandleFunc("/file-requests/{id}/send", SendFile).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
