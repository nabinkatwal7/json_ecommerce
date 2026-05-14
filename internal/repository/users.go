package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListUsers() ([]models.User, error) {
	return database.ReadJSON[models.User](s.usersFile())
}

func (s *Store) SaveUsers(users []models.User) error {
	return database.WriteJSON(s.usersFile(), users)
}

func (s *Store) FindUserByEmail(email string) (*models.User, error) {
	users, err := s.ListUsers()
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].Email == email {
			return &users[i], nil
		}
	}
	return nil, nil
}

func (s *Store) FindUserByID(id string) (*models.User, error) {
	users, err := s.ListUsers()
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].ID == id {
			return &users[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertUser(u models.User) error {
	users, err := s.ListUsers()
	if err != nil {
		return err
	}
	found := false
	for i := range users {
		if users[i].ID == u.ID {
			users[i] = u
			found = true
			break
		}
	}
	if !found {
		users = append(users, u)
	}
	return s.SaveUsers(users)
}
