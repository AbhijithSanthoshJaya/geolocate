package main

import (
	"log"
	"net/http"
	"time"

	"github.com/geolocate/server"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/getplace/{placeID}", server.GetPlacebyId).Methods("GET")
	r.HandleFunc("/geocode", server.GetGeocode).Methods("GET")
	r.HandleFunc("/geodecode", server.GetGeodecode).Methods("GET")

	// Start server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
