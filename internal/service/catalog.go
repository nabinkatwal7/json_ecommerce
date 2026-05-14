package service

import (
	"fmt"
	"strings"
	"time"

	"go-ecommerce-json/internal/cache"
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type CatalogService struct {
	Store    *repository.Store
	Cache    cache.CatalogCache
	CacheTTL time.Duration
}

func (c *CatalogService) cacheKeyProducts(categoryID, tagID string) string {
	return fmt.Sprintf("catalog:products:active:%s:%s", categoryID, tagID)
}

func (c *CatalogService) catalogTTL() time.Duration {
	if c.CacheTTL > 0 {
		return c.CacheTTL
	}
	return 30 * time.Second
}

func (c *CatalogService) bustCatalogCache() {
	if c.Cache != nil {
		c.Cache.InvalidateCatalog()
	}
}

func (c *CatalogService) ListActiveProducts(categoryID, tagID string) ([]models.Product, error) {
	tagID = strings.TrimSpace(tagID)
	if c.Cache != nil {
		var cached []models.Product
		if c.Cache.GetJSON(c.cacheKeyProducts(categoryID, tagID), &cached) {
			return cached, nil
		}
	}
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
		if tagID != "" && !productHasTagID(p, tagID) {
			continue
		}
		out = append(out, p)
	}
	if c.Cache != nil {
		c.Cache.SetJSON(c.cacheKeyProducts(categoryID, tagID), out, c.catalogTTL())
	}
	return out, nil
}

func productHasTagID(p models.Product, tagID string) bool {
	tagID = strings.TrimSpace(tagID)
	if tagID == "" {
		return true
	}
	for _, id := range p.TagIDs {
		if strings.TrimSpace(id) == tagID {
			return true
		}
	}
	return false
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

func (c *CatalogService) GetCategory(id string) (*models.Category, error) {
	cat, err := c.Store.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if cat == nil || !cat.IsActive {
		return nil, ErrNotFound
	}
	return cat, nil
}

func (c *CatalogService) AdminListProducts() ([]models.Product, error) {
	return c.Store.ListProducts()
}

// --- Admin ---

type ProductInput struct {
	Name        string                  `json:"name"`
	Slug        string                  `json:"slug"`
	Description string                  `json:"description"`
	Image       string                  `json:"image"`
	CategoryID  string                  `json:"categoryId"`
	Tags        []string                `json:"tags"`
	TagIDs      []string                `json:"tagIds"`
	Variants    []models.ProductVariant `json:"variants"`
	IsActive    bool                    `json:"isActive"`
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
	variants, err := normalizeVariants(in.Variants)
	if err != nil {
		return nil, err
	}
	if err := c.validateVariantSKUs("", variants); err != nil {
		return nil, err
	}
	if err := c.validateTagIDs(in.TagIDs); err != nil {
		return nil, err
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
		TagIDs:      append([]string(nil), in.TagIDs...),
		Variants:    variants,
		IsActive:    in.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := c.Store.UpsertProduct(p); err != nil {
		return nil, err
	}
	c.bustCatalogCache()
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
	variants, err := normalizeVariants(in.Variants)
	if err != nil {
		return nil, err
	}
	if err := c.validateVariantSKUs(p.ID, variants); err != nil {
		return nil, err
	}
	if err := c.validateTagIDs(in.TagIDs); err != nil {
		return nil, err
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
	p.TagIDs = append([]string(nil), in.TagIDs...)
	p.Variants = variants
	p.IsActive = in.IsActive
	p.UpdatedAt = now
	if err := c.Store.UpsertProduct(*p); err != nil {
		return nil, err
	}
	c.bustCatalogCache()
	return p, nil
}

func (c *CatalogService) AdminDeleteProduct(id string) error {
	if p, _ := c.Store.FindProductByID(id); p == nil {
		return ErrNotFound
	}
	if err := c.Store.DeleteProduct(id); err != nil {
		return err
	}
	c.bustCatalogCache()
	return nil
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
	c.bustCatalogCache()
	return &cat, nil
}

func (c *CatalogService) AdminUpdateCategory(id string, in CategoryInput) (*models.Category, error) {
	cat, err := c.Store.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrNotFound
	}
	slug := strings.TrimSpace(in.Slug)
	if slug != "" && slug != cat.Slug {
		cats, err := c.Store.ListCategories()
		if err != nil {
			return nil, err
		}
		for i := range cats {
			if cats[i].ID != id && cats[i].Slug == slug {
				return nil, ErrConflict
			}
		}
	}
	if strings.TrimSpace(in.Name) != "" {
		cat.Name = strings.TrimSpace(in.Name)
	}
	if slug != "" {
		cat.Slug = slug
	}
	cat.Description = in.Description
	cat.IsActive = in.IsActive
	if err := c.Store.UpsertCategory(*cat); err != nil {
		return nil, err
	}
	c.bustCatalogCache()
	return cat, nil
}

func (c *CatalogService) AdminDeleteCategory(id string) error {
	if cat, _ := c.Store.FindCategoryByID(id); cat == nil {
		return ErrNotFound
	}
	if err := c.Store.DeleteCategory(id); err != nil {
		return err
	}
	c.bustCatalogCache()
	return nil
}

func normalizeVariants(in []models.ProductVariant) ([]models.ProductVariant, error) {
	if len(in) == 0 {
		return nil, ErrValidation
	}
	out := make([]models.ProductVariant, len(in))
	copy(out, in)
	for i := range out {
		if strings.TrimSpace(out[i].SKU) == "" {
			return nil, ErrValidation
		}
		if strings.TrimSpace(out[i].ID) == "" {
			out[i].ID = uuid.NewString()
		}
	}
	seen := map[string]struct{}{}
	for _, v := range out {
		sku := strings.TrimSpace(v.SKU)
		if _, ok := seen[sku]; ok {
			return nil, ErrConflict
		}
		seen[sku] = struct{}{}
	}
	return out, nil
}

func (c *CatalogService) validateTagIDs(ids []string) error {
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			return ErrValidation
		}
		t, err := c.Store.FindTagByID(id)
		if err != nil {
			return err
		}
		if t == nil {
			return ErrValidation
		}
	}
	return nil
}

func (c *CatalogService) validateVariantSKUs(excludeProductID string, variants []models.ProductVariant) error {
	products, err := c.Store.ListProducts()
	if err != nil {
		return err
	}
	for _, nv := range variants {
		ns := strings.TrimSpace(nv.SKU)
		for _, p := range products {
			if p.ID == excludeProductID {
				continue
			}
			for _, v := range p.Variants {
				if strings.TrimSpace(v.SKU) == ns {
					return ErrConflict
				}
			}
		}
	}
	return nil
}
