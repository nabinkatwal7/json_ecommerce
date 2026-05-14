package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListProducts() ([]models.Product, error) {
	return database.ReadJSON[models.Product](s.productsFile())
}

func (s *Store) SaveProducts(products []models.Product) error {
	return database.WriteJSON(s.productsFile(), products)
}

func (s *Store) FindProductByID(id string) (*models.Product, error) {
	products, err := s.ListProducts()
	if err != nil {
		return nil, err
	}
	for i := range products {
		if products[i].ID == id {
			return &products[i], nil
		}
	}
	return nil, nil
}

func (s *Store) FindProductBySlug(slug string) (*models.Product, error) {
	products, err := s.ListProducts()
	if err != nil {
		return nil, err
	}
	for i := range products {
		if products[i].Slug == slug {
			return &products[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertProduct(p models.Product) error {
	products, err := s.ListProducts()
	if err != nil {
		return err
	}
	found := false
	for i := range products {
		if products[i].ID == p.ID {
			products[i] = p
			found = true
			break
		}
	}
	if !found {
		products = append(products, p)
	}
	return s.SaveProducts(products)
}

func (s *Store) DeleteProduct(id string) error {
	products, err := s.ListProducts()
	if err != nil {
		return err
	}
	out := products[:0]
	for _, p := range products {
		if p.ID != id {
			out = append(out, p)
		}
	}
	return s.SaveProducts(out)
}
