package types

// AppInfo represents a MicroAppInfo response
type AppInfo struct {
	Error  string `json:"error"`
	Common struct {
		StoreTags []string `json:"store_tags"`
	} `json:"common"`
}
