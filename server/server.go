package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/geolocate/client"
	"github.com/geolocate/geo"
	"github.com/gorilla/mux"
)

var defaultFieldMask = []geo.PlaceFieldMask{geo.PlaceFieldMaskBusinessStatus, geo.PlaceFieldMaskFormattedAddress, geo.PlaceFieldMaskDispName, geo.PlaceFieldMaskPlaceID, geo.PlaceFieldMaskTypes, geo.PlaceFieldMaskOpeningHours}
var apiKey = os.Getenv("API_KEY")
var header = geo.PlacesHeader{PlaceFieldMasks: defaultFieldMask, ApiKey: apiKey, ContentType: "application/json", MaskPrefix: false}

func GetPlacebyId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placeID := vars["placeID"]
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
		return
	}
	apiClient := geo.GeoClient{c}
	ctx := context.Background()
	place, err := apiClient.PlaceDetails(ctx, placeID, &header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((place))
}
func GetGeodecode(w http.ResponseWriter, r *http.Request) {
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
		return
	}
	queryParams := r.URL.Query()
	lat, err := strconv.ParseFloat(queryParams.Get("lat"), 64)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
	}
	long, _ := strconv.ParseFloat(queryParams.Get("long"), 64)
	testGeoClient := geo.GeoClient{c}
	ctx := context.Background()
	req := geo.GeocodingRequest{LatLng: &geo.LatLng{Lat: lat, Lng: long}}
	fmt.Printf("%+v/n", req)

	geocode, err := testGeoClient.Geodecode(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((geocode))
}
func GetGeocode(w http.ResponseWriter, r *http.Request) {
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
		return
	}
	queryParams := r.URL.Query()
	placeAddress := queryParams.Get("address")
	testGeoClient := geo.GeoClient{c}
	ctx := context.Background()
	req := geo.GeocodingRequest{Address: placeAddress}
	// fmt.Printf("%+v/n", req)
	geocode, err := testGeoClient.Geocode(ctx, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((geocode))
}
