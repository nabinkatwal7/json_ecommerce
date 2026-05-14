package service

import (
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"
)

// UserListsService persists wishlist and save-for-later per user (JSON store).
type UserListsService struct {
	Store *repository.Store
}

func (s *UserListsService) getOrEmptyWishlist(userID string) (models.Wishlist, error) {
	w, err := s.Store.FindWishlistByUser(userID)
	if err != nil {
		return models.Wishlist{}, err
	}
	if w == nil {
		return models.Wishlist{UserID: userID, Items: []models.WishlistItem{}}, nil
	}
	return *w, nil
}

func (s *UserListsService) ListWishlist(userID string) ([]models.WishlistItem, error) {
	w, err := s.getOrEmptyWishlist(userID)
	if err != nil {
		return nil, err
	}
	return w.Items, nil
}

func (s *UserListsService) AddWishlist(userID, productID, variantID string) ([]models.WishlistItem, error) {
	p, err := s.Store.FindProductByID(productID)
	if err != nil {
		return nil, err
	}
	if p == nil || !p.IsActive {
		return nil, ErrNotFound
	}
	var v *models.ProductVariant
	for i := range p.Variants {
		if p.Variants[i].ID == variantID {
			v = &p.Variants[i]
			break
		}
	}
	if v == nil {
		return nil, ErrValidation
	}
	w, err := s.getOrEmptyWishlist(userID)
	if err != nil {
		return nil, err
	}
	for _, it := range w.Items {
		if it.ProductID == productID && it.VariantID == variantID {
			return w.Items, nil
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	w.Items = append(w.Items, models.WishlistItem{
		ProductID: productID,
		VariantID: variantID,
		SKU:       v.SKU,
		Name:      p.Name,
		Price:     v.Price,
		Image:     p.Image,
		CreatedAt: now,
	})
	w.UpdatedAt = now
	if err := s.Store.UpsertWishlist(w); err != nil {
		return nil, err
	}
	return w.Items, nil
}

func (s *UserListsService) RemoveWishlist(userID, productID, variantID string) ([]models.WishlistItem, error) {
	w, err := s.getOrEmptyWishlist(userID)
	if err != nil {
		return nil, err
	}
	out := w.Items[:0]
	removed := false
	for _, it := range w.Items {
		if it.ProductID == productID && it.VariantID == variantID {
			removed = true
			continue
		}
		out = append(out, it)
	}
	if !removed {
		return nil, ErrNotFound
	}
	w.Items = out
	w.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.Store.UpsertWishlist(w); err != nil {
		return nil, err
	}
	return w.Items, nil
}

func (s *UserListsService) getOrEmptySaveLater(userID string) (models.SaveLaterList, error) {
	sl, err := s.Store.FindSaveLaterByUser(userID)
	if err != nil {
		return models.SaveLaterList{}, err
	}
	if sl == nil {
		return models.SaveLaterList{UserID: userID, Items: []models.SaveLaterItem{}}, nil
	}
	return *sl, nil
}

func (s *UserListsService) ListSaveLater(userID string) ([]models.SaveLaterItem, error) {
	sl, err := s.getOrEmptySaveLater(userID)
	if err != nil {
		return nil, err
	}
	return sl.Items, nil
}

func (s *UserListsService) AddSaveLater(userID, productID, variantID string) ([]models.SaveLaterItem, error) {
	p, err := s.Store.FindProductByID(productID)
	if err != nil {
		return nil, err
	}
	if p == nil || !p.IsActive {
		return nil, ErrNotFound
	}
	var v *models.ProductVariant
	for i := range p.Variants {
		if p.Variants[i].ID == variantID {
			v = &p.Variants[i]
			break
		}
	}
	if v == nil {
		return nil, ErrValidation
	}
	sl, err := s.getOrEmptySaveLater(userID)
	if err != nil {
		return nil, err
	}
	for _, it := range sl.Items {
		if it.ProductID == productID && it.VariantID == variantID {
			return sl.Items, nil
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	sl.Items = append(sl.Items, models.SaveLaterItem{
		ProductID: productID,
		VariantID: variantID,
		SKU:       v.SKU,
		Name:      p.Name,
		Price:     v.Price,
		Image:     p.Image,
		CreatedAt: now,
	})
	sl.UpdatedAt = now
	if err := s.Store.UpsertSaveLater(sl); err != nil {
		return nil, err
	}
	return sl.Items, nil
}

func (s *UserListsService) RemoveSaveLater(userID, productID, variantID string) ([]models.SaveLaterItem, error) {
	sl, err := s.getOrEmptySaveLater(userID)
	if err != nil {
		return nil, err
	}
	out := sl.Items[:0]
	removed := false
	for _, it := range sl.Items {
		if it.ProductID == productID && it.VariantID == variantID {
			removed = true
			continue
		}
		out = append(out, it)
	}
	if !removed {
		return nil, ErrNotFound
	}
	sl.Items = out
	sl.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.Store.UpsertSaveLater(sl); err != nil {
		return nil, err
	}
	return sl.Items, nil
}

// MoveToSaveLater removes from wishlist and adds to save-for-later (idempotent add).
func (s *UserListsService) MoveToSaveLater(userID, productID, variantID string) (wish []models.WishlistItem, later []models.SaveLaterItem, err error) {
	if _, err = s.RemoveWishlist(userID, productID, variantID); err != nil && err != ErrNotFound {
		return nil, nil, err
	}
	later, err = s.AddSaveLater(userID, productID, variantID)
	if err != nil {
		return nil, nil, err
	}
	wish, err = s.ListWishlist(userID)
	return wish, later, err
}

// MoveToWishlist removes from save-for-later and adds to wishlist.
func (s *UserListsService) MoveToWishlist(userID, productID, variantID string) (wish []models.WishlistItem, later []models.SaveLaterItem, err error) {
	if _, err = s.RemoveSaveLater(userID, productID, variantID); err != nil && err != ErrNotFound {
		return nil, nil, err
	}
	wish, err = s.AddWishlist(userID, productID, variantID)
	if err != nil {
		return nil, nil, err
	}
	later, err = s.ListSaveLater(userID)
	return wish, later, err
}
