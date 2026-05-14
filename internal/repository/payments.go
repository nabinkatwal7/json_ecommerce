package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListPayments() ([]models.Payment, error) {
	return database.ReadJSON[models.Payment](s.paymentsFile())
}

func (s *Store) SavePayments(payments []models.Payment) error {
	return database.WriteJSON(s.paymentsFile(), payments)
}

func (s *Store) ListPaymentsByOrderID(orderID string) ([]models.Payment, error) {
	all, err := s.ListPayments()
	if err != nil {
		return nil, err
	}
	var out []models.Payment
	for i := range all {
		if all[i].OrderID == orderID {
			out = append(out, all[i])
		}
	}
	return out, nil
}

func (s *Store) UpsertPayment(p models.Payment) error {
	payments, err := s.ListPayments()
	if err != nil {
		return err
	}
	found := false
	for i := range payments {
		if payments[i].ID == p.ID {
			payments[i] = p
			found = true
			break
		}
	}
	if !found {
		payments = append(payments, p)
	}
	return s.SavePayments(payments)
}
