package models

type Discount struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	Type          string  `json:"type"` // percent | fixed | bogo
	Value         float64 `json:"value"`
	MinimumAmount float64 `json:"minimumAmount"`
	IsActive      bool    `json:"isActive"`
	ExpiresAt     string  `json:"expiresAt"`
	// BOGO: for every (BuyQty + GetQty) units on eligible lines, discount GetQty units at line price (classic buy-one-get-one when 1/1).
	BuyQty     int    `json:"buyQty"`
	GetQty     int    `json:"getQty"`
	ProductID  string `json:"productId,omitempty"`
	CategoryID string `json:"categoryId,omitempty"`
}