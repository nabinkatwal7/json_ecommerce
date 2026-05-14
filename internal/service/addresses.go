package service

import (
	"strings"

	"go-ecommerce-json/internal/models"

	"github.com/google/uuid"
)

func (s *UserService) ListAddresses(userID string) ([]models.Address, error) {
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	return u.Addresses, nil
}

func (s *UserService) AddAddress(userID string, a models.Address) (*models.User, error) {
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	if strings.TrimSpace(a.FullName) == "" || strings.TrimSpace(a.AddressLine) == "" {
		return nil, ErrValidation
	}
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	if a.IsDefault {
		for i := range u.Addresses {
			u.Addresses[i].IsDefault = false
		}
	}
	u.Addresses = append(u.Addresses, a)
	if err := s.Store.UpsertUser(*u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) UpdateAddress(userID, addressID string, a models.Address) (*models.User, error) {
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	idx := -1
	for i := range u.Addresses {
		if u.Addresses[i].ID == addressID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, ErrNotFound
	}
	if strings.TrimSpace(a.FullName) == "" || strings.TrimSpace(a.AddressLine) == "" {
		return nil, ErrValidation
	}
	a.ID = addressID
	if a.IsDefault {
		for i := range u.Addresses {
			u.Addresses[i].IsDefault = false
		}
	}
	u.Addresses[idx] = a
	if err := s.Store.UpsertUser(*u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) DeleteAddress(userID, addressID string) (*models.User, error) {
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	out := u.Addresses[:0]
	removed := false
	for _, ad := range u.Addresses {
		if ad.ID == addressID {
			removed = true
			continue
		}
		out = append(out, ad)
	}
	if !removed {
		return nil, ErrNotFound
	}
	u.Addresses = out
	if err := s.Store.UpsertUser(*u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}
