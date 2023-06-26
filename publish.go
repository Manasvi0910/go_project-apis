package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Define the File struct
type File struct {
	ID      string  `json:"id"`
	Hash    string  `json:"hash"`
	Price   float64 `json:"price"`
	IsSaved bool    `json:"is_saved"`
}

var filesMutex sync.Mutex
var files = make(map[string]File)

// Define the API endpoint
func publishFile(w http.ResponseWriter, r *http.Request) {
	var file File
	_ = json.NewDecoder(r.Body).Decode(&file)

	// Check if the file already exists
	filesMutex.Lock()
	existingFile, exists := files[file.Hash]
	if exists {
		// If the file exists, update the price
		existingFile.Price = file.Price
		files[file.Hash] = existingFile
	} else {
		// If the file doesn't exist, create a new entry
		file.ID = generateUniqueID() // Generate a unique ID for the file
		file.IsSaved = true          // Simulating the saving of the file
		files[file.Hash] = file
	}
	filesMutex.Unlock()

	// Return the published file as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

// Helper function to generate a unique ID
func generateUniqueID() string {
	// Implement your logic to generate a unique ID here
	return "unique-id"
}

// main function to set up the API routes
func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/files", publishFile).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
