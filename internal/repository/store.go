package repository

import (
	"path/filepath"
)

// Store is the JSON file–backed persistence layer for all collections.
type Store struct {
	dir string
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name)
}

func (s *Store) usersFile() string       { return s.path("users.json") }
func (s *Store) productsFile() string    { return s.path("products.json") }
func (s *Store) categoriesFile() string  { return s.path("categories.json") }
func (s *Store) cartsFile() string        { return s.path("carts.json") }
func (s *Store) ordersFile() string       { return s.path("orders.json") }
func (s *Store) discountsFile() string  { return s.path("discounts.json") }
func (s *Store) paymentsFile() string     { return s.path("payments.json") }
func (s *Store) tagsFile() string         { return s.path("tags.json") }
func (s *Store) passwordResetsFile() string {
	return s.path("password_resets.json")
}
