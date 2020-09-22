package types

type FanAlgoliaOutgoing struct {
	Params string `json:"params"`
	APIKey string `json:"apiKey"`
}

type FanAlgoliaIncoming struct {
	Hits []struct {
		ProductID       string             `json:"product_id"`
		SKU             string             `json:"sku"`
		Name            string             `json:"name"`
		Slug            string             `json:"slug"`
		FullPrice       map[string]float32 `json:"fullPrice"`
		Presale         bool               `json:"presale"`
		Price           map[string]float32 `json:"price"`
		DiscountPercent uint               `json:"discount_percent"`
	}
}
