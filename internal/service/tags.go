package service

import (
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type TagService struct {
	Store *repository.Store
}

func (t *TagService) List() ([]models.Tag, error) {
	return t.Store.ListTags()
}

type TagInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (t *TagService) AdminCreate(in TagInput) (*models.Tag, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, ErrValidation
	}
	slug := strings.TrimSpace(in.Slug)
	if dup, _ := t.Store.FindTagBySlug(slug); dup != nil {
		return nil, ErrConflict
	}
	now := time.Now().UTC().Format(time.RFC3339)
	tag := models.Tag{
		ID:        uuid.NewString(),
		Name:      strings.TrimSpace(in.Name),
		Slug:      slug,
		CreatedAt: now,
	}
	if err := t.Store.UpsertTag(tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (t *TagService) AdminUpdate(id string, in TagInput) (*models.Tag, error) {
	tag, err := t.Store.FindTagByID(id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrNotFound
	}
	slug := strings.TrimSpace(in.Slug)
	if slug != "" && slug != tag.Slug {
		if dup, _ := t.Store.FindTagBySlug(slug); dup != nil {
			return nil, ErrConflict
		}
	}
	if strings.TrimSpace(in.Name) != "" {
		tag.Name = strings.TrimSpace(in.Name)
	}
	if slug != "" {
		tag.Slug = slug
	}
	if err := t.Store.UpsertTag(*tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *TagService) AdminDelete(id string) error {
	if tag, _ := t.Store.FindTagByID(id); tag == nil {
		return ErrNotFound
	}
	return t.Store.DeleteTag(id)
}
