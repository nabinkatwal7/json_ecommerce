package models

type Discount struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	Type          string  `json:"type"`
	Value         float64 `json:"value"`
	MinimumAmount float64 `json:"minimumAmount"`
	IsActive      bool    `json:"isActive"`
	ExpiresAt     string  `json:"expiresAt"`
}