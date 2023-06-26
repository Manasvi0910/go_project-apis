package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Order represents an order entity
type Order struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Quantity  int    `json:"quantity"`
	Fulfilled bool   `json:"fulfilled"`
}

var (
	ordersMutex sync.Mutex
	orders      = make(map[string]Order)
)

// CreateOrder creates a new order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)

	ordersMutex.Lock()
	// Check if order already exists
	if existingOrder, ok := orders[order.ID]; ok {
		// Order already exists, return the existing order
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingOrder)
		ordersMutex.Unlock()
		return
	}

	// Generate a unique ID for the order
	// You can implement your own unique ID generation logic here
	order.ID = "order-" + generateUniqueID()

	// Save the order
	orders[order.ID] = order
	ordersMutex.Unlock()

	// Return the created order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// FulfillOrder fulfills an existing order
func FulfillOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]

	ordersMutex.Lock()
	defer ordersMutex.Unlock()

	// Check if the order exists
	if order, ok := orders[orderID]; ok {
		// Check if the order is already fulfilled
		if order.Fulfilled {
			// Order is already fulfilled, return the existing order
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(order)
			return
		}

		// Update the order as fulfilled
		order.Fulfilled = true
		orders[orderID] = order

		// Return the fulfilled order
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	} else {
		// Order not found
		w.WriteHeader(http.StatusNotFound)
	}
}

// generateUniqueID generates a unique ID for the order
// You can replace this with your own unique ID generation logic
func generateUniqueID() string {
	// Implement your unique ID generation logic here
	return "unique-id"
}

func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/orders", CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}/fulfill", FulfillOrder).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
