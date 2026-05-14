package models

type Address struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Phone       string `json:"phone"`
	Country     string `json:"country"`
	State       string `json:"state"`
	City        string `json:"city"`
	PostalCode  string `json:"postalCode"`
	AddressLine string `json:"addressLine"`
	IsDefault   bool   `json:"isDefault"`
}

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Role         string    `json:"role"` // admin, customer
	Addresses    []Address `json:"addresses"`
	Segments     []string  `json:"segments,omitempty"` // derived marketing segments (refreshed periodically)
	CreatedAt    string    `json:"createdAt"`
}