package models

type CartItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	VariantID string  `json:"variantId"`
	Name      string  `json:"name"`
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Image     string  `json:"image"`
}

type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Items     []CartItem `json:"items"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}