package geo

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/geolocate/client"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Testing PlacesNearby encoder method. It takes a LatLng and Radius to find all establishments that match includeType
func Test_PlacesNearby(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	apiKey := os.Getenv("API_KEY")
	testclient, err := client.NewClient(client.AddAPIKey(apiKey))
	assert.NoError(t, err)
	assert.NotNil(t, testclient)
	testGeoClient := GeoClient{testclient}
	ctx := context.Background()
	incTypes := []PlaceType{"restaurant"}
	location := LocationRestriction{
		Circle{Center: Location{Latitude: 44.67775, Longitude: -63.67206}, Radius: 10000}}
	req := NearbySearchRequest{LocationRestriction: &location, MaxResultCount: 1, IncludedTypes: incTypes}
	fieldMask := []PlaceFieldMask{PlaceFieldMaskBusinessStatus, PlaceFieldMaskFormattedAddress, PlaceFieldMaskDispName, PlaceFieldMaskPlaceID, PlaceFieldMaskTypes, PlaceFieldMaskOpeningHours}
	header := PlacesHeader{FieldMasks: fieldMask, FieldMaskPrefix: true, TokenMask: ""}
	resp, err := testGeoClient.NearbySearch(ctx, &req, &header)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Testing PlacesNearby encoder method. It takes a LatLng and Radius to find all establishments that match text query.
//Eg bowling arena withing the locationBias expressed as LocationRestriction object Lat,Lng and Radius of search

func Test_TextSearch_locationBias(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	apiKey := os.Getenv("API_KEY")
	testclient, err := client.NewClient(client.AddAPIKey(apiKey), client.WithRateLimit(10))
	assert.NoError(t, err)
	assert.NotNil(t, testclient)
	testGeoClient := GeoClient{testclient}
	ctx := context.Background()
	textQuery := "bowling arena"
	locationBias := LocationRestriction{
		Circle{Center: Location{Latitude: 44.67775, Longitude: -63.67206}, Radius: 5000}}
	req := TextSearchRequest{TextQuery: textQuery, LocationBias: &locationBias, RankPreference: RankPreferenceDistance, PageSize: 5}
	fieldMask := []PlaceFieldMask{PlaceFieldMaskBusinessStatus, PlaceFieldMaskFormattedAddress, PlaceFieldMaskDispName, PlaceFieldMaskPlaceID, PlaceFieldMaskTypes, PlaceFieldMaskOpeningHours}
	header := PlacesHeader{FieldMasks: fieldMask, FieldMaskPrefix: true, TokenMask: ""}
	resp, err := testGeoClient.TextSearch(ctx, &req, &header)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_TextSearch_Restriction(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	apiKey := os.Getenv("API_KEY")
	testclient, err := client.NewClient(client.AddAPIKey(apiKey), client.WithRateLimit(10))
	assert.NoError(t, err)
	assert.NotNil(t, testclient)
	testGeoClient := GeoClient{testclient}
	ctx := context.Background()
	textQuery := "bowling arena"
	locationRestriction := RectangularRestriction{Rectangle: Rectangle{Low: Location{Latitude: 44.711211, Longitude: -63.722595}, High: Location{Latitude: 44.581167099, Longitude: -63.5431991}}}

	req := TextSearchRequest{TextQuery: textQuery, LocationRestriction: &locationRestriction, PageSize: 5}
	fieldMask := []PlaceFieldMask{PlaceFieldMaskBusinessStatus, PlaceFieldMaskFormattedAddress, PlaceFieldMaskDispName, PlaceFieldMaskPlaceID, PlaceFieldMaskTypes, PlaceFieldMaskOpeningHours}
	header := PlacesHeader{FieldMasks: fieldMask, FieldMaskPrefix: true, TokenMask: ""}
	resp, err := testGeoClient.TextSearch(ctx, &req, &header)
	assert.Error(t, err) // TODO: Issue giving viewport as locationrestriction. HTTP 400
	assert.Nil(t, resp)
}

func Test_PlaceDetails(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	apiKey := os.Getenv("API_KEY")
	testclient, err := client.NewClient(client.AddAPIKey(apiKey))
	assert.NoError(t, err)
	assert.NotNil(t, testclient)
	testGeoClient := GeoClient{testclient}
	ctx := context.Background()
	placeID := "ChIJy3Cb7veIWUsRDRRJADIvnms" // a real world location's placeID as set by Google
	fieldMask := []PlaceFieldMask{PlaceFieldMaskBusinessStatus, PlaceFieldMaskFormattedAddress, PlaceFieldMaskDispName, PlaceFieldMaskPlaceID, PlaceFieldMaskTypes, PlaceFieldMaskOpeningHours}
	header := PlacesHeader{FieldMasks: fieldMask, FieldMaskPrefix: false, TokenMask: ""}
	resp, err := testGeoClient.PlaceDetails(ctx, placeID, &header)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
