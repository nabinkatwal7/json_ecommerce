package models

type OrderItem struct {
	ProductID string  `json:"productId"`
	VariantID string  `json:"variantId"`
	Name      string  `json:"name"`
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type Order struct {
	ID              string      `json:"id"`
	UserID          string      `json:"userId"`
	Items           []OrderItem `json:"items"`
	ShippingAddress Address     `json:"shippingAddress"`
	Subtotal        float64     `json:"subtotal"`
	Discount        float64     `json:"discount"`
	Shipping        float64     `json:"shipping"`
	ShippingCarrier string      `json:"shippingCarrier,omitempty"` // e.g. fedex_ground, ups_ground, dhl_express, flat
	Total           float64     `json:"total"`
	Status          string      `json:"status"`        // created, paid, fulfilled, shipped, cancelled
	PaymentStatus   string      `json:"paymentStatus"` // pending, paid, failed, refunded
	InvoiceNumber   string      `json:"invoiceNumber,omitempty"`
	CreatedAt       string      `json:"createdAt"`
	UpdatedAt       string      `json:"updatedAt"`
}