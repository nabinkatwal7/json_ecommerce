package models

// PasswordResetToken is persisted so reset links survive restarts (JSON store).
type PasswordResetToken struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	TokenHash string `json:"tokenHash"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
}
