package repository

import (
	"strings"

	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListDiscounts() ([]models.Discount, error) {
	return database.ReadJSON[models.Discount](s.discountsFile())
}

func (s *Store) SaveDiscounts(discounts []models.Discount) error {
	return database.WriteJSON(s.discountsFile(), discounts)
}

func (s *Store) FindDiscountByCode(code string) (*models.Discount, error) {
	discounts, err := s.ListDiscounts()
	if err != nil {
		return nil, err
	}
	code = strings.TrimSpace(strings.ToUpper(code))
	for i := range discounts {
		if strings.EqualFold(strings.TrimSpace(discounts[i].Code), code) {
			return &discounts[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertDiscount(d models.Discount) error {
	discounts, err := s.ListDiscounts()
	if err != nil {
		return err
	}
	found := false
	for i := range discounts {
		if discounts[i].ID == d.ID {
			discounts[i] = d
			found = true
			break
		}
	}
	if !found {
		discounts = append(discounts, d)
	}
	return s.SaveDiscounts(discounts)
}
