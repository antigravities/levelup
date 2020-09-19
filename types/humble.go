package types

// HumbleSearch represents a Humble search query
type HumbleSearch struct {
	NumResults uint   `json:"num_results"`
	Query      string `json:"search"`
	Pages      uint   `json:"pages"`
	Results    []struct {
		HumanURL  string `json:"human_url"`
		FullPrice struct {
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount"`
		} `json:"full_price"`
		CurrentPrice struct {
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount"`
		} `json:"current_price"`
		Name string `json:"human_name"`
	} `json:"results"`
}
