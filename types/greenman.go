package types

// GreenManSearch represents a Green Man Gaming search query
type GreenManSearch struct {
	Results []struct {
		Name     string
		URL      string `json:"Url"`
		Editions []string
	} `json:"results"`
}
