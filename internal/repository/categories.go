package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListCategories() ([]models.Category, error) {
	return database.ReadJSON[models.Category](s.categoriesFile())
}

func (s *Store) SaveCategories(categories []models.Category) error {
	return database.WriteJSON(s.categoriesFile(), categories)
}

func (s *Store) FindCategoryByID(id string) (*models.Category, error) {
	categories, err := s.ListCategories()
	if err != nil {
		return nil, err
	}
	for i := range categories {
		if categories[i].ID == id {
			return &categories[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertCategory(c models.Category) error {
	categories, err := s.ListCategories()
	if err != nil {
		return err
	}
	found := false
	for i := range categories {
		if categories[i].ID == c.ID {
			categories[i] = c
			found = true
			break
		}
	}
	if !found {
		categories = append(categories, c)
	}
	return s.SaveCategories(categories)
}

func (s *Store) DeleteCategory(id string) error {
	categories, err := s.ListCategories()
	if err != nil {
		return err
	}
	out := categories[:0]
	for _, c := range categories {
		if c.ID != id {
			out = append(out, c)
		}
	}
	return s.SaveCategories(out)
}
