package types

// SteamStorefront represents the Storefront API response
type SteamStorefront struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string   `json:"type"`
		Name                string   `json:"name"`
		AppID               uint     `json:"steam_appid"`
		RequiredAge         uint     `json:"required_age"`
		IsFree              bool     `json:"is_free"`
		DetailedDescription string   `json:"detailed_description"`
		ShortDescription    string   `json:"short_description"`
		SupportedLanguages  string   `json:"supported_languages"`
		Reviews             string   `json:"Reviews"`
		HeaderImage         string   `json:"header_image"`
		Website             string   `json:"website"`
		LegalNotice         string   `json:"legal_notice"`
		Developers          []string `json:"developers"`
		Publishers          []string `json:"publishers"`
		Price               struct {
			Currency        string `json:"currency"`
			Initial         uint64 `json:"initial"`
			Final           uint64 `json:"final"`
			DiscountPercent uint   `json:"discount_percent"`
			Formatted       string `json:"final_formatted"`
		} `json:"price_overview"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
	} `json:"data"`
}

// SteamAppList represents the ISteamApps response
type SteamAppList struct {
	AppList struct {
		Apps []SteamGame `json:"apps"`
	} `json:"applist"`
}

// SteamGame represents an item in ISteamApps.Apps
type SteamGame struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
