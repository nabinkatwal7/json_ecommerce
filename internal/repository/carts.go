package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListCarts() ([]models.Cart, error) {
	return database.ReadJSON[models.Cart](s.cartsFile())
}

func (s *Store) SaveCarts(carts []models.Cart) error {
	return database.WriteJSON(s.cartsFile(), carts)
}

func (s *Store) FindCartByUserID(userID string) (*models.Cart, error) {
	carts, err := s.ListCarts()
	if err != nil {
		return nil, err
	}
	for i := range carts {
		if carts[i].UserID == userID {
			return &carts[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertCart(c models.Cart) error {
	carts, err := s.ListCarts()
	if err != nil {
		return err
	}
	found := false
	for i := range carts {
		if carts[i].ID == c.ID {
			carts[i] = c
			found = true
			break
		}
	}
	if !found {
		carts = append(carts, c)
	}
	return s.SaveCarts(carts)
}
