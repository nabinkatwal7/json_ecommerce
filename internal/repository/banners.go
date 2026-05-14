package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListBanners() ([]models.Banner, error) {
	return database.ReadJSON[models.Banner](s.bannersFile())
}

func (s *Store) SaveBanners(banners []models.Banner) error {
	return database.WriteJSON(s.bannersFile(), banners)
}

func (s *Store) FindBannerByID(id string) (*models.Banner, error) {
	list, err := s.ListBanners()
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == id {
			return &list[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertBanner(b models.Banner) error {
	list, err := s.ListBanners()
	if err != nil {
		return err
	}
	found := false
	for i := range list {
		if list[i].ID == b.ID {
			list[i] = b
			found = true
			break
		}
	}
	if !found {
		list = append(list, b)
	}
	return s.SaveBanners(list)
}

func (s *Store) DeleteBanner(id string) error {
	list, err := s.ListBanners()
	if err != nil {
		return err
	}
	out := list[:0]
	for _, b := range list {
		if b.ID == id {
			continue
		}
		out = append(out, b)
	}
	if len(out) == len(list) {
		return nil
	}
	return s.SaveBanners(out)
}
