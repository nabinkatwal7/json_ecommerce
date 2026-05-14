package models

type Payment struct {
	ID                string  `json:"id"`
	OrderID           string  `json:"orderId"`
	Provider          string  `json:"provider"` // stripe, paypal, etc.
	ProviderReference string  `json:"providerReference"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"` // pending, paid, failed, refunded
	CreatedAt         string  `json:"createdAt"`
}