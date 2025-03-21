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

// Testing GeoCode encoder method for a real world address that is known
func Test_Geocode(t *testing.T) {
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
	req := GeocodingRequest{Address: "29 Beechwood Terr,Halifax, Canada"}
	resp, err := testGeoClient.Geocode(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Testing ReverseGeoCode encoder method for a real world address that is known
func Test_ReverseGeocode(t *testing.T) {
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
	req := GeocodingRequest{LatLng: &LatLng{Lat: float64(44.67775), Lng: float64(-63.67206)}}
	resp, err := testGeoClient.ReverseGeocode(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
