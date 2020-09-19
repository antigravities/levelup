package types

// App represents a Steam app
type App struct {
	AppID         int        `json:"appid"`
	Name          string     `json:"name"`
	Prices        *AppPrices `json:"prices"`
	RecommendedBy []string   `json:"recommendedBy"`
	Developers    []string   `json:"developers"`
	Publishers    []string   `json:"publishers"`
	LastUpdate    int64      `json:"lastUpdate"`
	Description   string     `json:"description"`
	Genres        []string   `json:"genres"`
	ReleaseYear   int        `json:"release_year"`
}

// AppPrices represents the prices of an app
type AppPrices struct {
	Steam  map[string]AppPrice
	Humble map[string]AppPrice
}

// AppPrice represents the price of an app
type AppPrice struct {
	Price         uint64 `json:"price"`
	OriginalPrice uint64 `json:"original_price"`
	Discount      uint   `json:"discount"`
	URL           string `json:"url"`
}

var (
	// Regions is a list of regions supported by Level Up
	Regions []string = []string{"us", "fr"}
)
