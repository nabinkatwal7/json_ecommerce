package service

import (
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type PromoService struct {
	Store *repository.Store
}

type DiscountInput struct {
	Code          string  `json:"code"`
	Type          string  `json:"type"` // percent | fixed
	Value         float64 `json:"value"`
	MinimumAmount float64 `json:"minimumAmount"`
	IsActive      bool    `json:"isActive"`
	ExpiresAt     string  `json:"expiresAt"` // RFC3339 or empty
}

func normalizeDiscountType(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

func (p *PromoService) AdminCreateDiscount(in DiscountInput) (*models.Discount, error) {
	code := strings.TrimSpace(strings.ToUpper(in.Code))
	if code == "" {
		return nil, ErrValidation
	}
	t := normalizeDiscountType(in.Type)
	if t != "percent" && t != "fixed" {
		return nil, ErrValidation
	}
	if in.Value <= 0 {
		return nil, ErrValidation
	}
	if t == "percent" && in.Value > 100 {
		return nil, ErrValidation
	}
	if existing, _ := p.Store.FindDiscountByCode(code); existing != nil {
		return nil, ErrConflict
	}
	d := models.Discount{
		ID:            uuid.NewString(),
		Code:          code,
		Type:          t,
		Value:         in.Value,
		MinimumAmount: in.MinimumAmount,
		IsActive:      in.IsActive,
		ExpiresAt:     strings.TrimSpace(in.ExpiresAt),
	}
	if err := p.Store.UpsertDiscount(d); err != nil {
		return nil, err
	}
	return &d, nil
}

// ApplyDiscount returns the monetary discount to subtract from subtotal (>= 0).
func ApplyDiscount(d *models.Discount, subtotal float64, now time.Time) (float64, error) {
	if d == nil || !d.IsActive {
		return 0, ErrInactive
	}
	if d.ExpiresAt != "" {
		exp, err := time.Parse(time.RFC3339, d.ExpiresAt)
		if err != nil {
			return 0, ErrValidation
		}
		if now.After(exp) {
			return 0, ErrInactive
		}
	}
	if subtotal < d.MinimumAmount {
		return 0, ErrValidation
	}
	switch d.Type {
	case "percent":
		amt := subtotal * (d.Value / 100)
		if amt < 0 || amt > subtotal {
			return subtotal, nil
		}
		return amt, nil
	case "fixed":
		if d.Value >= subtotal {
			return subtotal, nil
		}
		return d.Value, nil
	default:
		return 0, ErrValidation
	}
}
