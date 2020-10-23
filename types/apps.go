package types

// App represents a Steam app
type App struct {
	AppID           int
	Name            string
	Prices          *AppPrices
	RecommendedAt   int64
	Developers      []string
	Publishers      []string
	LastUpdate      int64
	Description     string
	Genres          []string
	ReleaseYear     int
	Screenshot      string
	IsPending       bool
	ReviewsPositive int
	ReviewsTotal    int
	Demo            bool
	Score           float64
	Review          string
	Platforms       struct {
		Windows bool
		MacOS   bool
		Linux   bool
	}
}

// AppPrices represents the prices of an app
type AppPrices struct {
	Steam      map[string]AppPrice
	Humble     map[string]AppPrice
	Fanatical  map[string]AppPrice
	GameBillet map[string]AppPrice
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
	Regions []string = []string{"us", "fr", "uk"}
)
