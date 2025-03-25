package server

import (
	"context"
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

// Lookup a placeId to get all details of the place
func GetPlacebyId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placeID := vars["placeID"]
	if placeID == "" {
		responseJson(w, http.StatusBadRequest, nil)
		return
	}
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	apiClient := geo.GeoClient{c}
	ctx := context.Background()
	place, err := apiClient.PlaceDetails(ctx, placeID, &header)
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	responseJson(w, http.StatusOK, Response{Data: place, Error: ""}) // Success
}

// Look up  Geocoded Map input with lat,long and fetch a human readable address metadata

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
		responseJson(w, http.StatusBadRequest, err.Error())
		return
	}
	long, err := strconv.ParseFloat(queryParams.Get("long"), 64)
	if err != nil {
		responseJson(w, http.StatusBadRequest, err.Error())
		return
	}
	testGeoClient := geo.GeoClient{c}
	ctx := context.Background()
	req := geo.GeocodingRequest{LatLng: &geo.LatLng{Lat: lat, Lng: long}}
	geodecode, err := testGeoClient.Geodecode(ctx, &req)
	if err != nil {
		responseJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJson(w, http.StatusOK, Response{Data: geodecode, Error: ""}) // Success

}

// Look up a human readable address to get Geocoded Map response with lat,long and other geometric detail
func GetGeocode(w http.ResponseWriter, r *http.Request) {
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusInternalServerError, err.Error())
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
		responseJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJson(w, http.StatusOK, Response{Data: geocode, Error: ""}) // Success

}
