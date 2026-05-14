package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListWishlists() ([]models.Wishlist, error) {
	return database.ReadJSON[models.Wishlist](s.wishlistsFile())
}

func (s *Store) SaveWishlists(rows []models.Wishlist) error {
	return database.WriteJSON(s.wishlistsFile(), rows)
}

func (s *Store) FindWishlistByUser(userID string) (*models.Wishlist, error) {
	rows, err := s.ListWishlists()
	if err != nil {
		return nil, err
	}
	for i := range rows {
		if rows[i].UserID == userID {
			return &rows[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertWishlist(w models.Wishlist) error {
	rows, err := s.ListWishlists()
	if err != nil {
		return err
	}
	found := false
	for i := range rows {
		if rows[i].UserID == w.UserID {
			rows[i] = w
			found = true
			break
		}
	}
	if !found {
		rows = append(rows, w)
	}
	return s.SaveWishlists(rows)
}

func (s *Store) ListSaveLater() ([]models.SaveLaterList, error) {
	return database.ReadJSON[models.SaveLaterList](s.saveLaterFile())
}

func (s *Store) SaveSaveLater(rows []models.SaveLaterList) error {
	return database.WriteJSON(s.saveLaterFile(), rows)
}

func (s *Store) FindSaveLaterByUser(userID string) (*models.SaveLaterList, error) {
	rows, err := s.ListSaveLater()
	if err != nil {
		return nil, err
	}
	for i := range rows {
		if rows[i].UserID == userID {
			return &rows[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertSaveLater(w models.SaveLaterList) error {
	rows, err := s.ListSaveLater()
	if err != nil {
		return err
	}
	found := false
	for i := range rows {
		if rows[i].UserID == w.UserID {
			rows[i] = w
			found = true
			break
		}
	}
	if !found {
		rows = append(rows, w)
	}
	return s.SaveSaveLater(rows)
}
