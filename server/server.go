package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/geolocate/client"
	"github.com/geolocate/geo"
	"github.com/gorilla/mux"
)

var apiKey = os.Getenv("API_KEY")
var defaultFieldMask = []geo.PlaceFieldMask{geo.PlaceFieldMaskBusinessStatus, geo.PlaceFieldMaskFormattedAddress, geo.PlaceFieldMaskDispName, geo.PlaceFieldMaskPlaceID, geo.PlaceFieldMaskTypes, geo.PlaceFieldMaskOpeningHours}
var defaultPlaceHeader = geo.PlaceHeader{PlaceFieldMasks: defaultFieldMask, ApiKey: apiKey, ContentType: "application/json"}
var defaultIncType = geo.GetDefaultPlacesTypes() // If caller did not specify a valid place type, we go with default supported ones
var resultCount = int32(5)

// Look up  Geocoded Map input with lat,long and fetch a human readable address metadata

func GetGeodecode(w http.ResponseWriter, r *http.Request) {
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Error: err.Error()})

	}
	queryParams := r.URL.Query()
	lat, err := strconv.ParseFloat(queryParams.Get("latitude"), 64)
	if err != nil {
		responseJson(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	long, err := strconv.ParseFloat(queryParams.Get("longitude"), 64)
	if err != nil {
		responseJson(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	testGeoClient := geo.GeoClient{c}
	ctx := context.Background()
	req := geo.GeocodingRequest{LatLng: &geo.LatLng{Lat: lat, Lng: long}}
	geodecode, err := testGeoClient.Geodecode(ctx, &req)
	if err != nil {
		responseJson(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	responseJson(w, http.StatusOK, Response{Data: geodecode, Error: ""}) // Success

}

// Look up a human readable address to get Geocoded Map response with lat,long and other geometric detail
func GetGeocode(w http.ResponseWriter, r *http.Request) {
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Error: err.Error()})
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
		responseJson(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	responseJson(w, http.StatusOK, Response{Data: geocode}) // Success

}

// Define a struct to match the expected JSON body
type PlacesNearby struct {
	Lat    float64  `json:"latitude"`
	Long   float64  `json:"longitude"`
	Radius int64    `json:"radius"`
	Types  []string `json:"types"`
}

// Define a struct to match the expected JSON body
type PlacesFromText struct {
	Lat       float64 `json:"latitude"`
	Long      float64 `json:"longitude"`
	Radius    int64   `json:"radius"`
	Text      string  `json:"text,omitempty"`
	PageToken string  `json:"pageToken,omitempty"`
}

// Find Places Nearby a user. Filter out places using incTypes to get results that match user preferences
func GetPlacesNearby(w http.ResponseWriter, r *http.Request) {
	var params PlacesNearby
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		responseJson(w, http.StatusBadRequest, Response{Data: nil, Error: err.Error()})
		return
	}
	incTypes := defaultIncType
	if len(params.Types) > 0 { // if we got a valid Type from user. Need to lookup Type A and B table to find match.TODO
		var placeTypes []geo.PlaceType
		for _, t := range params.Types {
			placeTypes = append(placeTypes, geo.PlaceType(t))
		}
		incTypes = placeTypes // If place types are invalid, then we will have a problem. Need to think this further.TODO
	}
	location := geo.LocationRestriction{Circle: geo.Circle{Center: geo.Location{Latitude: params.Lat, Longitude: params.Long}, Radius: params.Radius}}
	req := geo.NearbySearchRequest{LocationRestriction: &location, MaxResultCount: resultCount, IncludedTypes: incTypes}
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	apiClient := geo.GeoClient{c}
	ctx := context.Background()
	header := geo.PlacesHeader{PlaceHeader: defaultPlaceHeader, MaskPrefix: true}
	place, err := apiClient.NearbySearch(ctx, &req, &header) //TODO
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	responseJson(w, http.StatusOK, Response{Data: place, Error: ""}) // Success
}

// Find Places Nearby a user. Filter out places using incTypes to get results that match user preferences
func GetPlacesFromText(w http.ResponseWriter, r *http.Request) {
	var params PlacesFromText
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		responseJson(w, http.StatusBadRequest, Response{Data: nil, Error: err.Error()})
		return
	}
	if params.Text == "" {
		responseJson(w, http.StatusBadRequest, Response{Data: nil, Error: "Please enter a valid search text"})
		return
	}
	textQuery := params.Text
	locationBias := geo.LocationRestriction{Circle: geo.Circle{Center: geo.Location{Latitude: params.Lat, Longitude: params.Long}, Radius: params.Radius}}
	req := geo.TextSearchRequest{TextQuery: textQuery, LocationBias: &locationBias, RankPreference: geo.RankPreferenceDistance, PageSize: resultCount, PageToken: params.PageToken}
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	apiClient := geo.GeoClient{c}
	ctx := context.Background()
	fieldMask := make([]geo.PlaceFieldMask, len(defaultFieldMask))
	copy(fieldMask, defaultFieldMask)
	fieldMask = append(fieldMask, geo.PlaceFieldNextPageToken)
	placeHeader := geo.PlaceHeader{PlaceFieldMasks: fieldMask, ApiKey: apiKey, ContentType: "application/json"}

	header := geo.PlacesHeader{PlaceHeader: placeHeader, MaskPrefix: true}
	place, err := apiClient.TextSearch(ctx, &req, &header) //TODO
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	responseJson(w, http.StatusOK, Response{Data: place, Error: ""}) // Success
}

// Lookup a placeId to get all details of the place
func GetPlacebyId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placeID := vars["placeID"]
	if placeID == "" {
		responseJson(w, http.StatusBadRequest, Response{Data: nil, Error: "Please enter a valid placeId"})
		return
	}
	c, err := client.NewClient(client.AddAPIKey(apiKey))
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	apiClient := geo.GeoClient{c}
	ctx := context.Background()
	header := geo.PlacesHeader{PlaceHeader: defaultPlaceHeader, MaskPrefix: false}
	place, err := apiClient.PlaceDetails(ctx, placeID, &header)
	if err != nil {
		responseJson(w, http.StatusServiceUnavailable, Response{Data: nil, Error: err.Error()})
		return
	}
	responseJson(w, http.StatusOK, Response{Data: place, Error: ""}) // Success
}
func GetAllTypes(w http.ResponseWriter, r *http.Request) {
	placeTypes := geo.GetAllPlacesTypes()
	responseJson(w, http.StatusOK, Response{Data: placeTypes, Error: ""}) // Success

}

func GetDefaultTypes(w http.ResponseWriter, r *http.Request) {
	placeTypes := geo.GetDefaultPlacesTypes()
	responseJson(w, http.StatusOK, Response{Data: placeTypes, Error: ""}) // Success

}
