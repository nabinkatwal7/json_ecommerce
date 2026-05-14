package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListOrders() ([]models.Order, error) {
	return database.ReadJSON[models.Order](s.ordersFile())
}

func (s *Store) SaveOrders(orders []models.Order) error {
	return database.WriteJSON(s.ordersFile(), orders)
}

func (s *Store) FindOrderByID(id string) (*models.Order, error) {
	orders, err := s.ListOrders()
	if err != nil {
		return nil, err
	}
	for i := range orders {
		if orders[i].ID == id {
			return &orders[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertOrder(o models.Order) error {
	orders, err := s.ListOrders()
	if err != nil {
		return err
	}
	found := false
	for i := range orders {
		if orders[i].ID == o.ID {
			orders[i] = o
			found = true
			break
		}
	}
	if !found {
		orders = append(orders, o)
	}
	return s.SaveOrders(orders)
}

func (s *Store) OrdersByUser(userID string) ([]models.Order, error) {
	orders, err := s.ListOrders()
	if err != nil {
		return nil, err
	}
	var out []models.Order
	for _, o := range orders {
		if o.UserID == userID {
			out = append(out, o)
		}
	}
	return out, nil
}
