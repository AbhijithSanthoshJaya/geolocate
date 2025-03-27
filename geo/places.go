package geo

import (
	"context"
	"errors"
	"os"
	"strings"

	// Included for image/jpeg's decoder
	_ "image/jpeg"

	"github.com/geolocate/client"
)

// New Places API endpoints that need to be queried for all Place related Search
type placesAPI struct {
	Host     string
	BasePath string
}

var places = placesAPI{Host: "https://places.googleapis.com", BasePath: "/v1/places"}

// Converts PlacesHeader into a map to be used as HTTP header in POST request to Places API
func (h *PlacesHeader) Headers() map[string]string {
	header := map[string]string{}
	prefix := ""
	if h.FieldMaskPrefix {
		prefix = "places." // Only for Places(plural) requests. For looking up a single place, we dont need this prefix. This api is wierd
	}
	fieldMaskHeader := FieldMaskHeader(h.FieldMasks, prefix, h.TokenMask)
	header["X-Goog-Api-Key"] = os.Getenv("API_KEY")
	header["X-Goog-FieldMask"] = strings.Join(fieldMaskHeader, ",")
	header["Content-Type"] = "application/json"
	return header
}

// API Response to call to Places API
type PlacesSearchResponse struct {
	Places        []Place `json:"places"`
	NextPageToken string  `json:"nextPageToken"`
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
	Location            Location       `json:"location"`
	BusinessStatus      BusinessStatus `json:"businessStatus"`
	PhoneNumber         string         `json:"nationalPhoneNumber"`
	Photos              []Photo        `json:"photos,omitempty"`
	Timezone            Timezone       `json:"timeZone,omitempty"`
	RegularOpeningHours OpeningHours   `json:"regularOpeningHours,omitempty"`
}

type NearbySearchRequest struct {
	RegionCode           string               `json:"regionCode,omitempty"`
	IncludedTypes        []PlaceType          `json:"includedTypes,omitempty"`
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
	PageSize                         int32                   `json:"pageSize,omitempty"`
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
	var response PlacesSearchResponse
	api := &client.ApiConfig{
		Host: places.Host,
		Path: places.BasePath + ":searchNearby",
	}
	if err := c.JsonPost(ctx, api, r, h, &response); err != nil {
		return PlacesSearchResponse{}, err
	}
	return response, nil
}

// TextSearch lets you search for places within a specified area that matches user text input. You can refine
// your search request by supplying the text and location restictions you are searching for.
func (c *GeoClient) TextSearch(ctx context.Context, r *TextSearchRequest, h *PlacesHeader) (PlacesSearchResponse, error) {
	if r.TextQuery == "" && r.PageToken == "" {
		return PlacesSearchResponse{}, errors.New("maps: Required fields Text Search and nextPage token are empty")
	}
	var response PlacesSearchResponse
	api := &client.ApiConfig{
		Host: places.Host,
		Path: places.BasePath + ":searchText",
	}
	if err := c.JsonPost(ctx, api, r, h, &response); err != nil {
		return PlacesSearchResponse{}, err
	}
	return response, nil

}

// Lookup a place details using placeID.
func (c *GeoClient) PlaceDetails(ctx context.Context, id string, h *PlacesHeader) (Place, error) {

	var response Place
	api := &client.ApiConfig{
		Host: places.Host,
		Path: places.BasePath + "/" + id,
	}
	if err := c.JsonGet(ctx, api, nil, h, &response); err != nil {
		return Place{}, err
	}
	return response, nil
}
func GetAllPlacesTypes() []PlaceType {
	AllPlaceTypes := []PlaceType{
		AcaiShop, AfghaniRestaurant, AfricanRestaurant, AmericanRestaurant,
		AsianRestaurant, BagelShop, Bakery, Bar, BarAndGrill, BarbecueRestaurant,
		BrazilianRestaurant, BreakfastRestaurant, BrunchRestaurant, BuffetRestaurant,
		Cafe, Cafeteria, CandyStore, CatCafe, ChineseRestaurant, ChocolateFactory,
		ChocolateShop, CoffeeShop, Confectionery, Deli, DessertRestaurant,
		DessertShop, Diner, DogCafe, DonutShop, FastFoodRestaurant,
		FineDiningRestaurant, FoodCourt, FrenchRestaurant, GreekRestaurant,
		HamburgerRestaurant, IceCreamShop, IndianRestaurant, IndonesianRestaurant,
		ItalianRestaurant, JapaneseRestaurant, JuiceShop, KoreanRestaurant,
		LebaneseRestaurant, MealDelivery, MealTakeaway, MediterraneanRestaurant,
		MexicanRestaurant, MiddleEasternRestaurant, PizzaRestaurant, Pub,
		RamenRestaurant, Restaurant, SandwichShop, SeafoodRestaurant,
		SpanishRestaurant, SteakHouse, SushiRestaurant, TeaHouse, ThaiRestaurant,
		TurkishRestaurant, VeganRestaurant, VegetarianRestaurant,
		VietnameseRestaurant, WineBar,
	}
	return AllPlaceTypes
}
func GetDefaultPlacesTypes() []PlaceType {
	DefaultPlaces := []PlaceType{
		AsianRestaurant, BagelShop, Bakery, Bar, BarAndGrill, BarbecueRestaurant, BreakfastRestaurant, BrunchRestaurant, BuffetRestaurant,
		Cafe, CatCafe, ChocolateShop, CoffeeShop, DessertRestaurant,
		DessertShop, Diner, DogCafe, FastFoodRestaurant,
		FineDiningRestaurant, IceCreamShop, IndianRestaurant, IndonesianRestaurant,
		ItalianRestaurant, JapaneseRestaurant, JuiceShop, KoreanRestaurant,
		LebaneseRestaurant, MediterraneanRestaurant,
		MexicanRestaurant, MiddleEasternRestaurant, PizzaRestaurant, Pub,
		RamenRestaurant, Restaurant, SeafoodRestaurant,
		SpanishRestaurant, SteakHouse, SushiRestaurant, TeaHouse, ThaiRestaurant,
		TurkishRestaurant, VeganRestaurant, VegetarianRestaurant,
		VietnameseRestaurant, WineBar,
	}
	return DefaultPlaces
}

type PlacesHeader struct {
	FieldMasks      []PlaceFieldMask
	FieldMaskPrefix bool
	TokenMask       string
}
type StaticHeader struct {
	ContentType string
	ApiKey      string
	FieldMasks  []PlaceFieldMask
}

func FieldMaskHeader(placeFieldMasks []PlaceFieldMask, prefix string, tokenMask string) []string {
	var fieldMask []string
	for _, field := range placeFieldMasks {
		fieldMask = append(fieldMask, string(prefix+string(field)))
	}
	if tokenMask != "" {
		fieldMask = append(fieldMask, tokenMask)
	}
	return fieldMask
}

type LocationRestriction struct {
	Circle Circle `json:"circle"`
}
type Circle struct {
	Center Location `json:"center"`
	Radius int64    `json:"radius"`
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
const MaskNextPageToken = "nextPageToken"
