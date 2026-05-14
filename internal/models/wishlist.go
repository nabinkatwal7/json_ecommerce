package models

type WishlistItem struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Price     float64 `json:"price"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdAt"`
}

type Wishlist struct {
	UserID    string         `json:"userId"`
	Items     []WishlistItem `json:"items"`
	UpdatedAt string         `json:"updatedAt"`
}

type SaveLaterItem struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Price     float64 `json:"price"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdAt"`
}

type SaveLaterList struct {
	UserID    string          `json:"userId"`
	Items     []SaveLaterItem `json:"items"`
	UpdatedAt string          `json:"updatedAt"`
}
