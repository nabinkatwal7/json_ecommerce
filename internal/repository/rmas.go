package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListRMAs() ([]models.RMA, error) {
	return database.ReadJSON[models.RMA](s.rmasFile())
}

func (s *Store) SaveRMAs(rows []models.RMA) error {
	return database.WriteJSON(s.rmasFile(), rows)
}

func (s *Store) UpsertRMA(r models.RMA) error {
	rows, err := s.ListRMAs()
	if err != nil {
		return err
	}
	found := false
	for i := range rows {
		if rows[i].ID == r.ID {
			rows[i] = r
			found = true
			break
		}
	}
	if !found {
		rows = append(rows, r)
	}
	return s.SaveRMAs(rows)
}

func (s *Store) FindRMAByID(id string) (*models.RMA, error) {
	rows, err := s.ListRMAs()
	if err != nil {
		return nil, err
	}
	for i := range rows {
		if rows[i].ID == id {
			return &rows[i], nil
		}
	}
	return nil, nil
}

func (s *Store) RMAsByUser(userID string) ([]models.RMA, error) {
	rows, err := s.ListRMAs()
	if err != nil {
		return nil, err
	}
	var out []models.RMA
	for _, r := range rows {
		if r.UserID == userID {
			out = append(out, r)
		}
	}
	return out, nil
}
