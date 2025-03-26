package geo

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

// PlaceType represents different types of food & beverage establishments
type PlaceType string

// List of PlaceTypes as constants
const (
	AcaiShop                PlaceType = "acai_shop"
	AfghaniRestaurant       PlaceType = "afghani_restaurant"
	AfricanRestaurant       PlaceType = "african_restaurant"
	AmericanRestaurant      PlaceType = "american_restaurant"
	AsianRestaurant         PlaceType = "asian_restaurant"
	BagelShop               PlaceType = "bagel_shop"
	Bakery                  PlaceType = "bakery"
	Bar                     PlaceType = "bar"
	BarAndGrill             PlaceType = "bar_and_grill"
	BarbecueRestaurant      PlaceType = "barbecue_restaurant"
	BrazilianRestaurant     PlaceType = "brazilian_restaurant"
	BreakfastRestaurant     PlaceType = "breakfast_restaurant"
	BrunchRestaurant        PlaceType = "brunch_restaurant"
	BuffetRestaurant        PlaceType = "buffet_restaurant"
	Cafe                    PlaceType = "cafe"
	Cafeteria               PlaceType = "cafeteria"
	CandyStore              PlaceType = "candy_store"
	CatCafe                 PlaceType = "cat_cafe"
	ChineseRestaurant       PlaceType = "chinese_restaurant"
	ChocolateFactory        PlaceType = "chocolate_factory"
	ChocolateShop           PlaceType = "chocolate_shop"
	CoffeeShop              PlaceType = "coffee_shop"
	Confectionery           PlaceType = "confectionery"
	Deli                    PlaceType = "deli"
	DessertRestaurant       PlaceType = "dessert_restaurant"
	DessertShop             PlaceType = "dessert_shop"
	Diner                   PlaceType = "diner"
	DogCafe                 PlaceType = "dog_cafe"
	DonutShop               PlaceType = "donut_shop"
	FastFoodRestaurant      PlaceType = "fast_food_restaurant"
	FineDiningRestaurant    PlaceType = "fine_dining_restaurant"
	FoodCourt               PlaceType = "food_court"
	FrenchRestaurant        PlaceType = "french_restaurant"
	GreekRestaurant         PlaceType = "greek_restaurant"
	HamburgerRestaurant     PlaceType = "hamburger_restaurant"
	IceCreamShop            PlaceType = "ice_cream_shop"
	IndianRestaurant        PlaceType = "indian_restaurant"
	IndonesianRestaurant    PlaceType = "indonesian_restaurant"
	ItalianRestaurant       PlaceType = "italian_restaurant"
	JapaneseRestaurant      PlaceType = "japanese_restaurant"
	JuiceShop               PlaceType = "juice_shop"
	KoreanRestaurant        PlaceType = "korean_restaurant"
	LebaneseRestaurant      PlaceType = "lebanese_restaurant"
	MealDelivery            PlaceType = "meal_delivery"
	MealTakeaway            PlaceType = "meal_takeaway"
	MediterraneanRestaurant PlaceType = "mediterranean_restaurant"
	MexicanRestaurant       PlaceType = "mexican_restaurant"
	MiddleEasternRestaurant PlaceType = "middle_eastern_restaurant"
	PizzaRestaurant         PlaceType = "pizza_restaurant"
	Pub                     PlaceType = "pub"
	RamenRestaurant         PlaceType = "ramen_restaurant"
	Restaurant              PlaceType = "restaurant"
	SandwichShop            PlaceType = "sandwich_shop"
	SeafoodRestaurant       PlaceType = "seafood_restaurant"
	SpanishRestaurant       PlaceType = "spanish_restaurant"
	SteakHouse              PlaceType = "steak_house"
	SushiRestaurant         PlaceType = "sushi_restaurant"
	TeaHouse                PlaceType = "tea_house"
	ThaiRestaurant          PlaceType = "thai_restaurant"
	TurkishRestaurant       PlaceType = "turkish_restaurant"
	VeganRestaurant         PlaceType = "vegan_restaurant"
	VegetarianRestaurant    PlaceType = "vegetarian_restaurant"
	VietnameseRestaurant    PlaceType = "vietnamese_restaurant"
	WineBar                 PlaceType = "wine_bar"
)
