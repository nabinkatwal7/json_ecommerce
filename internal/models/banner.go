package models

// Banner is storefront marketing content (carousel slides or top-bar announcements).
type Banner struct {
	ID        string `json:"id"`
	Slot      string `json:"slot"` // home_carousel | announcement
	Title     string `json:"title"`
	Body      string `json:"body"` // subtitle or announcement copy
	ImageURL  string `json:"imageUrl"`
	LinkURL   string `json:"linkUrl"`
	SortOrder int    `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
	StartsAt  string `json:"startsAt,omitempty"` // RFC3339; empty = no start bound
	EndsAt    string `json:"endsAt,omitempty"`   // RFC3339; empty = no end bound
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
