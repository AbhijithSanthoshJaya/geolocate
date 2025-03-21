package geo

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/geolocate/client"
)

var geocodingAPI = &client.ApiConfig{
	Host: "https://maps.googleapis.com",
	Path: "/maps/api/geocode/json",
}

// LatLng represents a location on the Earth.
type LatLng struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

// LatLngBounds represents a bounded square area on the Earth.
type LatLngBounds struct {
	NorthEast LatLng `json:"northeast"`
	SouthWest LatLng `json:"southwest"`
}

func (l *LatLng) String() string {
	return strconv.FormatFloat(l.Lat, 'f', -1, 64) +
		"," +
		strconv.FormatFloat(l.Lng, 'f', -1, 64)
}

type Component string

// GeocodeAccuracy is the type of a location result from the Geocoding API.
type GeocodeAccuracy string

const (
	//Use it to restrict results accurate down to street address precision.
	GeocodeAccuracyRooftop = GeocodeAccuracy("ROOFTOP")
	//Use it to restrict results to those that reflect an approximation
	// interpolated between two precise points.
	GeocodeAccuracyRI = GeocodeAccuracy("RANGE_INTERPOLATED")
	//Use it to restrict results to geometric centers of a location such as
	// a polyline or polygon.
	GeocodeAccuracyGC = GeocodeAccuracy("GEOMETRIC_CENTER")
	//Use it to restrict results to those that are characterized as approximate.
	GeocodeAccuracyApprox = GeocodeAccuracy("APPROXIMATE")
)

type GeoClient struct {
	*client.Client
}

// GeocodingRequest is the request structure for Geocoding API. It includes fields for both encoding and reverse geocoding
type GeocodingRequest struct {
	// Geocoding fields

	// Address is the street address that you want to geocode, in the format used by local post service
	Address string
	// Components is a component filter geocode generation. Please refer to original doc:
	// https://developers.google.com/maps/documentation/geocoding/intro#ComponentFiltering
	Components map[Component]string
	// Region is the region code, specified as a ccTLD two-character value. Optional.
	Region string
	// Reverse geocoding fields

	// Reverse Geocoding.
	// LocationType is an array of one or more geocoding accuracy types. Optional.
	LocationType []GeocodeAccuracy
	// LatLng is the textual latitude/longitude value for which you wish to obtain the
	// closest, human-readable address.
	LatLng *LatLng
	// ResultType is an array of one or more address types. Optional.
	ResultType []string
	// PlaceID is a string which contains the place_id, which can be used for reverse
	// geocoding requests. Either LatLng or PlaceID is required for Reverse Geocoding.
	PlaceID string
	// Language is the language in which to return results. Optional.
	Language string
	// Pass custom api query params to backend
	Custom url.Values
}

// GeocodingResponse is the response to a Geocoding API request.
type GeocodingResponse struct {
	// Results is the Geocoding results
	Results []GeocodingResult
}

// GeocodingResult is a single geocoded address
type GeocodingResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          AddressGeometry    `json:"geometry"`
	Types             []string           `json:"types"`
	PlaceID           string             `json:"place_id"`
	NavigationPoints  []Location         `json:"navigation_points"`

	// PartialMatch indicates that the geocoder did not return an exact match for
	// the original request, though it was able to match part of the requested address.
	// You may wish to examine the original request for misspellings and/or an incomplete address.
	PartialMatch bool `json:"partial_match"`
	// Plus codes can be used as a replacement for street addresses in places where they do not exist
	// where most buildings are not numbered or streets are not accurately named.
	// In our app, we only want to provide locations with valid address
	PlusCode PlusCode `json:"plus_code"`
}
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AddressWithCode (see https://en.wikipedia.org/wiki/Open_Location_Code and https://plus.codes/)
// is an encoded location reference, derived from latitude and longitude coordinates,
type PlusCode struct {
	// GlobalCode is a 4 character area code and 6 character or longer local code (849VCWC8+R9).
	GlobalCode string `json:"global_code"`
	// CompoundCode is a 6 character or longer local code with an explicit location (CWC8+R9, Mountain View, CA, USA).
	CompoundCode string `json:"compound_code"`
}

// AddressComponent is a part of an address
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// AddressGeometry is the location of a an address
type AddressGeometry struct {
	Location     LatLng       `json:"location"`
	LocationType string       `json:"location_type"`
	Bounds       LatLngBounds `json:"bounds"`
	Viewport     LatLngBounds `json:"viewport"`
	Types        []string     `json:"types"`
}

// Fetch Geo Feed from Google Maps API. Method can be used to convert an humanreadable address into Geocoded Result
func (c *GeoClient) Geocode(ctx context.Context, r *GeocodingRequest) (GeocodingResponse, error) {
	//Either address or component is necessary
	if r.Address == "" && len(r.Components) == 0 {
		return GeocodingResponse{}, errors.New("maps: Required fields address and/or components are all missing")
	}
	var response struct {
		Results []GeocodingResult `json:"results"`
		respStatus
	}

	if err := c.JsonGet(ctx, geocodingAPI, r, &response); err != nil {
		return GeocodingResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return GeocodingResponse{}, err
	}

	return GeocodingResponse{response.Results}, nil
}

// ReverseGeocode makes a Reverse Geocoding API request, returning a human readable address
// from Geocoded(Lat,Lng) or placeID type input
func (c *GeoClient) ReverseGeocode(ctx context.Context, r *GeocodingRequest) (GeocodingResponse, error) {
	// Either LatLng or PlaceID is necessary
	if r.Address != "" {
		return GeocodingResponse{}, errors.New("maps: Addr field must be empty,provide only latlng or placeID")

	}
	if r.LatLng == nil && r.PlaceID == "" {
		return GeocodingResponse{}, errors.New("maps: Required fields LatLng and/or PlaceID are both missing")
	}

	var response struct {
		Results []GeocodingResult `json:"results"`
		respStatus
	}

	if err := c.JsonGet(ctx, geocodingAPI, r, &response); err != nil {
		return GeocodingResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return GeocodingResponse{}, err
	}

	return GeocodingResponse{response.Results}, nil
}

// Converts a GeocodingRequest struct to a http queryparam object to pass to the http request
func (r *GeocodingRequest) Params() url.Values {
	q := make(url.Values)

	for k, v := range r.Custom {
		q[k] = v
	}
	if r.Address != "" {
		q.Set("address", r.Address)
	}
	var cf []string
	for c, f := range r.Components {
		cf = append(cf, string(c)+":"+f)
	}
	if len(cf) > 0 {
		q.Set("components", strings.Join(cf, "|"))
	}
	if r.Region != "" {
		q.Set("region", r.Region)
	}
	if r.LatLng != nil {
		q.Set("latlng", r.LatLng.String())
	}
	if len(r.ResultType) > 0 {
		q.Set("result_type", strings.Join(r.ResultType, "|"))
	}
	if len(r.LocationType) > 0 {
		var lt []string
		for _, l := range r.LocationType {
			lt = append(lt, string(l))
		}
		q.Set("location_type", strings.Join(lt, "|"))
	}
	if r.PlaceID != "" {
		q.Set("place_id", r.PlaceID)
	}
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	return q
}

// Contains the common response fields to most API calls inside
// the Google Maps APIs. This can be used to parse Status Response found
// at end of every Google Maps API response
type respStatus struct {
	// Status contains the status of the request, and may contain debugging
	// information to help you track down why the call failed.
	Status string `json:"status"`
	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// StatusError returns an error if this object has a Status different
// from OK or ZERO_RESULTS. While this happens, the http response is still success
//
//	witt HTTP_200
//
// This is different from API returning HTTP_ERRORS(5xx.4xx etc). Those
// are handled separately
func (c *respStatus) StatusError() error {
	if c.Status != "OK" && c.Status != "ZERO_RESULTS" {
		return fmt.Errorf("maps: %s - %s", c.Status, c.ErrorMessage)
	}
	return nil
}
