package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"go-ecommerce-json/internal/auth"
	"go-ecommerce-json/internal/mail"
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type PasswordResetService struct {
	Store    *repository.Store
	Mail     *mail.Sender
	AppURL   string // e.g. https://shop.example.com
	TokenTTL time.Duration
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// RequestPasswordReset always returns nil error for anti-enumeration when mail is configured.
func (p *PasswordResetService) RequestPasswordReset(in ForgotPasswordInput) error {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" {
		return ErrValidation
	}
	u, err := p.Store.FindUserByEmail(email)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return err
	}
	token := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()
	exp := now.Add(p.TokenTTL)
	row := models.PasswordResetToken{
		ID:        uuid.NewString(),
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: exp.Format(time.RFC3339),
		CreatedAt: now.Format(time.RFC3339),
	}
	if err := p.Store.UpsertPasswordReset(row); err != nil {
		return err
	}
	if !p.Mail.Configured() {
		return nil
	}
	link := strings.TrimRight(p.AppURL, "/") + "/reset-password?token=" + token
	if p.AppURL == "" {
		link = "(configure APP_PUBLIC_URL) token=" + token
	}
	body := "Reset your password using this link (expires soon):\n\n" + link + "\n"
	_ = p.Mail.SendPlain(u.Email, "Password reset", body)
	return nil
}

func (p *PasswordResetService) ResetPassword(in ResetPasswordInput) error {
	if len(in.NewPassword) < 8 || strings.TrimSpace(in.Token) == "" {
		return ErrValidation
	}
	sum := sha256.Sum256([]byte(strings.TrimSpace(in.Token)))
	hash := hex.EncodeToString(sum[:])
	rows, err := p.Store.ListPasswordResets()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	var match *models.PasswordResetToken
	for i := range rows {
		if rows[i].TokenHash != hash {
			continue
		}
		exp, err := time.Parse(time.RFC3339, rows[i].ExpiresAt)
		if err != nil {
			continue
		}
		if now.After(exp) {
			continue
		}
		match = &rows[i]
		break
	}
	if match == nil {
		return ErrUnauthorized
	}
	u, err := p.Store.FindUserByID(match.UserID)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrNotFound
	}
	newHash, err := auth.HashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = newHash
	if err := p.Store.UpsertUser(*u); err != nil {
		return err
	}
	_ = p.Store.DeletePasswordReset(match.ID)
	return nil
}
