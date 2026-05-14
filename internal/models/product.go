package models

type ProductVariant struct {
	ID    string  `json:"id"`
	SKU   string  `json:"sku"`
	Size  string  `json:"size"`
	Color string  `json:"color"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type Product struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	Image       string           `json:"image"`
	CategoryID  string           `json:"categoryId"`
	Tags        []string         `json:"tags"`
	Variants    []ProductVariant `json:"variants"`
	IsActive    bool             `json:"isActive"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}