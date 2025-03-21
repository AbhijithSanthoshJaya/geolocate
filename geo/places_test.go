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
	incTypes := []string{"restaurant"}
	location := LocationRestriction{
		Circle{Center: LatLng{Lat: 44.67775, Lng: -63.67206}, Radius: 10000}}
	req := NearbySearchRequest{LocationRestriction: &location, MaxResultCount: 1, IncludedTypes: incTypes}
	fieldMask := []PlaceFieldMask{PlaceFieldMaskBusinessStatus, PlaceFieldMaskFormattedAddress, PlaceFieldMaskDispName, PlaceFieldMaskPlaceID, PlaceFieldMaskTypes, PlaceFieldMaskOpeningHours}
	header := PlacesHeader{PlaceFieldMasks: fieldMask, ApiKey: apiKey, ContentType: "application/json"}
	resp, err := testGeoClient.NearbySearch(ctx, &req, &header)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
