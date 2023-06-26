package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Define the FileCreator struct
type FileCreator struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	FileFormat string `json:"file_format"`
}

// Define the Hoster struct
type Hoster struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Define the API endpoints
func createFileCreator(w http.ResponseWriter, r *http.Request) {
	var creator FileCreator
	_ = json.NewDecoder(r.Body).Decode(&creator)

	// Perform necessary operations to store the file creator details

	// Return the created file creator as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(creator)
}

func createHoster(w http.ResponseWriter, r *http.Request) {
	var hoster Hoster
	_ = json.NewDecoder(r.Body).Decode(&hoster)

	// Perform necessary operations to store the hoster details

	// Return the created hoster as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hoster)
}

func matchFileCreatorsWithHosters(w http.ResponseWriter, r *http.Request) {
	// Perform necessary operations to match file creators with hosters

	// Return the matched results as JSON response
	w.Header().Set("Content-Type", "application/json")
	// Prepare the matched results
	results := []struct {
		Creator FileCreator
		Hoster  Hoster
	}{
		{
			Creator: FileCreator{
				ID:         "1",
				Name:       "File Creator 1",
				Location:   "Location 1",
				FileFormat: "Format 1",
			},
			Hoster: Hoster{
				ID:       "A",
				Name:     "Hoster A",
				Location: "Location A",
			},
		},
		// Add more matched results here
	}
	json.NewEncoder(w).Encode(results)
}

// main function to set up the API routes
func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/file-creators", createFileCreator).Methods("POST")
	r.HandleFunc("/hosters", createHoster).Methods("POST")
	r.HandleFunc("/match", matchFileCreatorsWithHosters).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
