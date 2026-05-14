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
	Type          string  `json:"type"` // percent | fixed | bogo
	Value         float64 `json:"value"`
	MinimumAmount float64 `json:"minimumAmount"`
	IsActive      bool    `json:"isActive"`
	ExpiresAt     string  `json:"expiresAt"`
	BuyQty        int     `json:"buyQty"`
	GetQty        int     `json:"getQty"`
	ProductID     string  `json:"productId"`
	CategoryID    string  `json:"categoryId"`
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
	switch t {
	case "percent", "fixed":
		if in.Value <= 0 {
			return nil, ErrValidation
		}
		if t == "percent" && in.Value > 100 {
			return nil, ErrValidation
		}
	case "bogo":
		buy, get := in.BuyQty, in.GetQty
		if buy <= 0 {
			buy = 1
		}
		if get <= 0 {
			get = 1
		}
	default:
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
		BuyQty:        in.BuyQty,
		GetQty:        in.GetQty,
		ProductID:     strings.TrimSpace(in.ProductID),
		CategoryID:    strings.TrimSpace(in.CategoryID),
	}
	if d.Type == "bogo" {
		if d.BuyQty <= 0 {
			d.BuyQty = 1
		}
		if d.GetQty <= 0 {
			d.GetQty = 1
		}
	}
	if err := p.Store.UpsertDiscount(d); err != nil {
		return nil, err
	}
	return &d, nil
}

func validateDiscountWindow(d *models.Discount, subtotal float64, now time.Time) error {
	if d == nil || !d.IsActive {
		return ErrInactive
	}
	if d.ExpiresAt != "" {
		exp, err := time.Parse(time.RFC3339, d.ExpiresAt)
		if err != nil {
			return ErrValidation
		}
		if now.After(exp) {
			return ErrInactive
		}
	}
	if subtotal < d.MinimumAmount {
		return ErrValidation
	}
	return nil
}

// ApplyDiscount returns the monetary discount for percent/fixed types (>= 0).
func ApplyDiscount(d *models.Discount, subtotal float64, now time.Time) (float64, error) {
	if err := validateDiscountWindow(d, subtotal, now); err != nil {
		return 0, err
	}
	switch normalizeDiscountType(d.Type) {
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

func bogoLineEligible(store *repository.Store, d *models.Discount, line models.OrderItem) bool {
	p, err := store.FindProductByID(line.ProductID)
	if err != nil || p == nil {
		return false
	}
	if strings.TrimSpace(d.ProductID) != "" && p.ID != d.ProductID {
		return false
	}
	if strings.TrimSpace(d.CategoryID) != "" && p.CategoryID != d.CategoryID {
		return false
	}
	return true
}

// ApplyBogoDiscount computes BOGO savings per eligible line (buy N get M free units at the line unit price).
func ApplyBogoDiscount(store *repository.Store, d *models.Discount, items []models.OrderItem, subtotal float64, now time.Time) (float64, error) {
	if err := validateDiscountWindow(d, subtotal, now); err != nil {
		return 0, err
	}
	buy, get := d.BuyQty, d.GetQty
	if buy <= 0 {
		buy = 1
	}
	if get <= 0 {
		get = 1
	}
	bundle := buy + get
	var discount float64
	for _, line := range items {
		if !bogoLineEligible(store, d, line) {
			continue
		}
		bundles := line.Quantity / bundle
		freeUnits := bundles * get
		discount += float64(freeUnits) * line.Price
	}
	if discount > subtotal {
		discount = subtotal
	}
	if discount < 0 {
		discount = 0
	}
	return discount, nil
}

// ComputeDiscountAmount routes discount logic for checkout (percent, fixed, or BOGO).
func ComputeDiscountAmount(store *repository.Store, d *models.Discount, subtotal float64, items []models.OrderItem, now time.Time) (float64, error) {
	switch normalizeDiscountType(d.Type) {
	case "percent", "fixed":
		return ApplyDiscount(d, subtotal, now)
	case "bogo":
		return ApplyBogoDiscount(store, d, items, subtotal, now)
	default:
		return 0, ErrValidation
	}
}
