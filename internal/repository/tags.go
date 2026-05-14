package repository

import (
	"go-ecommerce-json/internal/database"
	"go-ecommerce-json/internal/models"
)

func (s *Store) ListTags() ([]models.Tag, error) {
	return database.ReadJSON[models.Tag](s.tagsFile())
}

func (s *Store) SaveTags(tags []models.Tag) error {
	return database.WriteJSON(s.tagsFile(), tags)
}

func (s *Store) FindTagByID(id string) (*models.Tag, error) {
	tags, err := s.ListTags()
	if err != nil {
		return nil, err
	}
	for i := range tags {
		if tags[i].ID == id {
			return &tags[i], nil
		}
	}
	return nil, nil
}

func (s *Store) FindTagBySlug(slug string) (*models.Tag, error) {
	tags, err := s.ListTags()
	if err != nil {
		return nil, err
	}
	for i := range tags {
		if tags[i].Slug == slug {
			return &tags[i], nil
		}
	}
	return nil, nil
}

func (s *Store) UpsertTag(t models.Tag) error {
	tags, err := s.ListTags()
	if err != nil {
		return err
	}
	found := false
	for i := range tags {
		if tags[i].ID == t.ID {
			tags[i] = t
			found = true
			break
		}
	}
	if !found {
		tags = append(tags, t)
	}
	return s.SaveTags(tags)
}

func (s *Store) DeleteTag(id string) error {
	tags, err := s.ListTags()
	if err != nil {
		return err
	}
	out := tags[:0]
	for _, t := range tags {
		if t.ID != id {
			out = append(out, t)
		}
	}
	return s.SaveTags(out)
}
