package geo

import (
	"context"
	"errors"
	"strings"

	// Included for image/jpeg's decoder
	_ "image/jpeg"

	"github.com/geolocate/client"
)

// New Places API endpoints that need to be queried for all Nearby Search
var placesSearchAPI = &client.ApiConfig{
	Host:  "https://places.googleapis.com",
	FPath: "/v1/places",
	Path:  "",
}

// Converts PlacesHeader into a map to be used as HTTP header in POST request to Places API
func (h *PlacesHeader) Headers() map[string]string {
	header := map[string]string{}
	prefix := ""
	if h.MaskPrefix {
		prefix = "places." // Only for Places(plural) requests. For looking up a single place, we dont need this prefix. This api is wierd
	}
	fieldMaskHeader := FieldMaskHeader(h.PlaceFieldMasks, prefix)
	header["X-Goog-Api-Key"] = h.ApiKey
	header["X-Goog-FieldMask"] = strings.Join(fieldMaskHeader, ",")
	header["Content-Type"] = h.ContentType
	return header
}

// API Response to call to Places API
type PlacesSearchResponse struct {
	Places []Place `json:"places"`
}
type LocalizedText struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

// Object representing a 'Place' as represented by Google Places
type Place struct {
	Id                  string         `json:"id"`
	DisplayName         LocalizedText  `json:"displayName"`
	Types               []string       `json:"types"`
	FormattedAddress    string         `json:"formattedAddress"`
	Rating              int32          `json:"rating"`
	Location            *LatLng        `json:"location"`
	BusinessStatus      BusinessStatus `json:"businessStatus"`
	PhoneNumber         string         `json:"nationalPhoneNumber"`
	Photos              []Photo        `json:"photos,omitempty"`
	Timezone            Timezone       `json:"timeZone,omitempty"`
	RegularOpeningHours OpeningHours   `json:"regularOpeningHours,omitempty"`
}

type NearbySearchRequest struct {
	RegionCode           string               `json:"regionCode,omitempty"`
	IncludedTypes        []string             `json:"includedTypes,omitempty"`
	ExcludedTypes        []string             `json:"excludedTypes,omitempty"`
	IncludedPrimaryTypes []string             `json:"includedPrimaryTypes,omitempty"`
	ExcludedPrimaryTypes []string             `json:"excludedPrimaryTypes,omitempty"`
	MaxResultCount       int32                `json:"maxResultCount,omitempty"`
	LocationRestriction  *LocationRestriction `json:"locationRestriction"`
	RankPreference       RankPreference       `json:"rankPreference,omitempty"`
}

type TextSearchRequest struct {
	TextQuery                        string                  `json:"textQuery"`
	IncludedType                     string                  `json:"includedType,omitempty"`
	IncludePureServiceAreaBusinesses bool                    `json:"includePureServiceAreaBusinesses,omitempty"`
	PageSize                         int                     `json:"pageSize,omitempty"`
	PageToken                        string                  `json:"pageToken,omitempty"`
	StrictTypeFiltering              bool                    `json:"strictTypeFiltering,omitempty"`
	LocationBias                     *LocationRestriction    `json:"locationBias,omitempty"`
	RankPreference                   RankPreference          `json:"rankPreference,omitempty"`
	LocationRestriction              *RectangularRestriction `json:"locationRestriction,omitempty"`
}

// NearbySearch lets you search for places within a specified area. You can refine
// your search request by supplying the type of place you are searching for.
func (c *GeoClient) NearbySearch(ctx context.Context, r *NearbySearchRequest, h *PlacesHeader) (PlacesSearchResponse, error) {
	if r.LocationRestriction == nil {
		return PlacesSearchResponse{}, errors.New("maps: Required fields LocationRestriction missing")
	}
	var response struct {
		Places []Place `json:"places"`
	}
	placesSearchAPI.Path = placesSearchAPI.FPath + ":searchNearby"
	if err := c.JsonPost(ctx, placesSearchAPI, r, h, &response); err != nil {
		return PlacesSearchResponse{}, err
	}
	return PlacesSearchResponse{response.Places}, nil
}

// TextSearch lets you search for places within a specified area that matches user text input. You can refine
// your search request by supplying the text and location restictions you are searching for.
func (c *GeoClient) TextSearch(ctx context.Context, r *TextSearchRequest, h *PlacesHeader) (PlacesSearchResponse, error) {
	if r.TextQuery == "" {
		return PlacesSearchResponse{}, errors.New("maps: Required fields Text Search is empty")
	}
	var response struct {
		Places []Place `json:"places"`
	}
	placesSearchAPI.Path = placesSearchAPI.FPath + ":searchText"

	if err := c.JsonPost(ctx, placesSearchAPI, r, h, &response); err != nil {
		return PlacesSearchResponse{}, err
	}
	return PlacesSearchResponse{response.Places}, nil

}

// Lookup a place details using placeID.
func (c *GeoClient) PlaceDetails(ctx context.Context, id string, h *PlacesHeader) (Place, error) {

	var response Place
	placesSearchAPI.Path = placesSearchAPI.FPath + "/" + id
	if err := c.JsonGet(ctx, placesSearchAPI, nil, h, &response); err != nil {
		return Place{}, err
	}
	return response, nil
}

type PlacesHeader struct {
	ContentType     string
	ApiKey          string
	PlaceFieldMasks []PlaceFieldMask
	MaskPrefix      bool
}

func FieldMaskHeader(placeFieldMasks []PlaceFieldMask, prefix string) []string {
	var fieldMask []string
	for _, fields := range placeFieldMasks {
		fieldMask = append(fieldMask, string(prefix+string(fields)))
	}
	return fieldMask
}

type LocationRestriction struct {
	Circle Circle `json:"circle"`
}
type Circle struct {
	Center LatLng `json:"center"`
	Radius int64  `json:"radius"`
}
type RectangularRestriction struct {
	Rectangle Rectangle `json:"rectangle"`
}
type Rectangle struct {
	Low  LatLng `json:"low"`
	High LatLng `json:"high"`
}
type RankPreference string

const (
	RankPreferenceUnspecified = RankPreference("RANK_PREFERENCE_UNSPECIFIED")
	RankPreferenceDistance    = RankPreference("DISTANCE")
	RankPreferencePopularity  = RankPreference("POPULARITY")
)

type PlaceTypes string

const (
	IncludedTypesRestaurant = PlaceTypes("restaurant")
	IncludedTypeBar         = PlaceTypes("bar")
)

type PlaceFieldMask string

// The individual Places Field Masks to trim fields in result returned by API.
const (
	PlaceFieldMaskDispName             = PlaceFieldMask("displayName")
	PlaceFieldMaskDineIn               = PlaceFieldMask("dineIn")
	PlaceFieldMaskFormattedAddress     = PlaceFieldMask("formattedAddress")
	PlaceFieldMaskFormattedPhoneNumber = PlaceFieldMask("nationalPhoneNumber")
	PlaceFieldMaskBusinessStatus       = PlaceFieldMask("businessStatus")
	PlaceFieldMaskPhotos               = PlaceFieldMask("photos")
	PlaceFieldMaskPlaceID              = PlaceFieldMask("id")
	PlaceFieldMaskRatings              = PlaceFieldMask("rating")
	PlaceFieldMaskTypes                = PlaceFieldMask("types")
	PlaceFieldMaskOpeningHours         = PlaceFieldMask("regularOpeningHours")
)
