package geo

// PlaceType restricts Place API search to the results to places matching the
// specified type.
type PlaceType string

// PriceLevel is the Price Levels for Places API
type PriceLevel string

// Price Levels for the Places API
const (
	PriceLevelFree          = PriceLevel("0")
	PriceLevelInexpensive   = PriceLevel("1")
	PriceLevelModerate      = PriceLevel("2")
	PriceLevelExpensive     = PriceLevel("3")
	PriceLevelVeryExpensive = PriceLevel("4")
)

type BusinessStatus string

const (
	BusinessStatusUnspecified       = BusinessStatus("BUSINESS_STATUS_UNSPECIFIED")
	BusinessStatusOperational       = BusinessStatus("OPERATIONAL")
	BusinessStatusClosedTemporarily = BusinessStatus("CLOSED_TEMPORARILY")
	BusinessStatusClosedPermanently = BusinessStatus("CLOSED_PERMANENTLY")
)

type OpeningHours struct {
	Periods             []Period
	WeekdayDescriptions []string
	SecondaryHoursType  SecondaryHoursType
	SpecialDays         []SpecialDay
	NextOpenTime        string
	NextCloseTime       string
	OpenNow             bool
}
type SecondaryHoursType string

type SpecialDay struct {
	Date Date `json:"date"`
}

// Photo describes a photo available with a Search Result.
type Photo struct {
	// PhotoReference is used to identify the photo when you perform a Photo request.
	Name string `json:"name"`
	// Height is the maximum height of the image.
	Height int `json:"heightPx"`
	// Width is the maximum width of the image.
	Width int `json:"widthPx"`
	// htmlAttributions contains any required attributions.
	FlagContentUri string `json:"flagContentUri"`
	GoogleMapsUri  string `json:"googleMapsUri"`
}
type Timezone struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type Period struct {
	Open  Point `json:"open"`
	Close Point `json:"close"`
}
type Point struct {
	Date      Date `json:"date"`
	Truncated bool `json:"truncated"`
	Day       int  `json:"day"`
	Hour      int  `json:"hour"`
	Minute    int  `json:"minute"`
}
type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}
