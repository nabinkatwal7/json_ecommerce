package models

type RMAItem struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Price     float64 `json:"price"`
}

// RMA is a return merchandise authorization tied to a shipped/paid order.
type RMA struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	OrderID     string    `json:"orderId"`
	Items       []RMAItem `json:"items"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"` // requested, approved, received, refunded, rejected
	AdminNote   string    `json:"adminNote,omitempty"`
	RefundAmount float64  `json:"refundAmount,omitempty"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}
