package service

import (
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type CatalogService struct {
	Store *repository.Store
}

func (c *CatalogService) ListActiveProducts(categoryID string) ([]models.Product, error) {
	all, err := c.Store.ListProducts()
	if err != nil {
		return nil, err
	}
	var out []models.Product
	for _, p := range all {
		if !p.IsActive {
			continue
		}
		if categoryID != "" && p.CategoryID != categoryID {
			continue
		}
		out = append(out, p)
	}
	return out, nil
}

func (c *CatalogService) GetProduct(id string) (*models.Product, error) {
	p, err := c.Store.FindProductByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil || !p.IsActive {
		return nil, ErrNotFound
	}
	return p, nil
}

func (c *CatalogService) GetProductBySlug(slug string) (*models.Product, error) {
	p, err := c.Store.FindProductBySlug(slug)
	if err != nil {
		return nil, err
	}
	if p == nil || !p.IsActive {
		return nil, ErrNotFound
	}
	return p, nil
}

func (c *CatalogService) ListActiveCategories() ([]models.Category, error) {
	all, err := c.Store.ListCategories()
	if err != nil {
		return nil, err
	}
	var out []models.Category
	for _, cat := range all {
		if cat.IsActive {
			out = append(out, cat)
		}
	}
	return out, nil
}

// --- Admin ---

type ProductInput struct {
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Image       string                 `json:"image"`
	CategoryID  string                 `json:"categoryId"`
	Tags        []string               `json:"tags"`
	Variants    []models.ProductVariant `json:"variants"`
	IsActive    bool                   `json:"isActive"`
}

func (c *CatalogService) AdminCreateProduct(in ProductInput) (*models.Product, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, ErrValidation
	}
	if dup, _ := c.Store.FindProductBySlug(strings.TrimSpace(in.Slug)); dup != nil {
		return nil, ErrConflict
	}
	if cat, _ := c.Store.FindCategoryByID(in.CategoryID); in.CategoryID != "" && cat == nil {
		return nil, ErrValidation
	}
	now := time.Now().UTC().Format(time.RFC3339)
	p := models.Product{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(in.Name),
		Slug:        strings.TrimSpace(in.Slug),
		Description: in.Description,
		Image:       in.Image,
		CategoryID:  in.CategoryID,
		Tags:        in.Tags,
		Variants:    in.Variants,
		IsActive:    in.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := c.Store.UpsertProduct(p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *CatalogService) AdminUpdateProduct(id string, in ProductInput) (*models.Product, error) {
	p, err := c.Store.FindProductByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrNotFound
	}
	slug := strings.TrimSpace(in.Slug)
	if slug != "" && slug != p.Slug {
		if dup, _ := c.Store.FindProductBySlug(slug); dup != nil {
			return nil, ErrConflict
		}
	}
	if in.CategoryID != "" {
		if cat, _ := c.Store.FindCategoryByID(in.CategoryID); cat == nil {
			return nil, ErrValidation
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if strings.TrimSpace(in.Name) != "" {
		p.Name = strings.TrimSpace(in.Name)
	}
	if slug != "" {
		p.Slug = slug
	}
	p.Description = in.Description
	p.Image = in.Image
	p.CategoryID = in.CategoryID
	p.Tags = in.Tags
	p.Variants = in.Variants
	p.IsActive = in.IsActive
	p.UpdatedAt = now
	if err := c.Store.UpsertProduct(*p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *CatalogService) AdminDeleteProduct(id string) error {
	if p, _ := c.Store.FindProductByID(id); p == nil {
		return ErrNotFound
	}
	return c.Store.DeleteProduct(id)
}

type CategoryInput struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

func (c *CatalogService) AdminCreateCategory(in CategoryInput) (*models.Category, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, ErrValidation
	}
	cats, err := c.Store.ListCategories()
	if err != nil {
		return nil, err
	}
	slug := strings.TrimSpace(in.Slug)
	for i := range cats {
		if cats[i].Slug == slug {
			return nil, ErrConflict
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	cat := models.Category{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(in.Name),
		Slug:        slug,
		Description: in.Description,
		IsActive:    in.IsActive,
		CreatedAt:   now,
	}
	if err := c.Store.UpsertCategory(cat); err != nil {
		return nil, err
	}
	return &cat, nil
}
