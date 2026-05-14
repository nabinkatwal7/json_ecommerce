package service

import (
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"
)

// adjustVariantStock changes inventory for each order line. sign=-1 deducts, sign=+1 restores.
func adjustVariantStock(store *repository.Store, items []models.OrderItem, sign int) error {
	if len(items) == 0 {
		return nil
	}
	products, err := store.ListProducts()
	if err != nil {
		return err
	}
	byID := make(map[string]int, len(products))
	for i := range products {
		byID[products[i].ID] = i
	}
	// Validate before mutating (avoid partial writes on error).
	for _, line := range items {
		idx, ok := byID[line.ProductID]
		if !ok {
			return ErrValidation
		}
		p := &products[idx]
		vidx := -1
		for i := range p.Variants {
			if p.Variants[i].ID == line.VariantID {
				vidx = i
				break
			}
		}
		if vidx < 0 {
			return ErrValidation
		}
		if sign < 0 && p.Variants[vidx].Stock < line.Quantity {
			return ErrInsufficientStock
		}
	}
	for _, line := range items {
		idx := byID[line.ProductID]
		p := &products[idx]
		vidx := -1
		for i := range p.Variants {
			if p.Variants[i].ID == line.VariantID {
				vidx = i
				break
			}
		}
		delta := line.Quantity * sign
		p.Variants[vidx].Stock += delta
	}
	return store.SaveProducts(products)
}
