package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListPasswordResets() ([]models.PasswordResetToken, error) {
	return database.ReadJSON[models.PasswordResetToken](s.passwordResetsFile())
}

func (s *Store) SavePasswordResets(rows []models.PasswordResetToken) error {
	return database.WriteJSON(s.passwordResetsFile(), rows)
}

func (s *Store) UpsertPasswordReset(row models.PasswordResetToken) error {
	rows, err := s.ListPasswordResets()
	if err != nil {
		return err
	}
	found := false
	for i := range rows {
		if rows[i].ID == row.ID {
			rows[i] = row
			found = true
			break
		}
	}
	if !found {
		rows = append(rows, row)
	}
	return s.SavePasswordResets(rows)
}

func (s *Store) DeletePasswordReset(id string) error {
	rows, err := s.ListPasswordResets()
	if err != nil {
		return err
	}
	out := rows[:0]
	for _, r := range rows {
		if r.ID != id {
			out = append(out, r)
		}
	}
	return s.SavePasswordResets(out)
}
