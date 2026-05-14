package service

import (
	"sort"
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type BannerService struct {
	Store *repository.Store
}

func normalizeBannerSlot(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "home_carousel", "announcement":
		return s
	default:
		return ""
	}
}

func bannerVisibleAt(b models.Banner, now time.Time) bool {
	if !b.IsActive {
		return false
	}
	if b.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, b.StartsAt)
		if err == nil && now.Before(t) {
			return false
		}
	}
	if b.EndsAt != "" {
		t, err := time.Parse(time.RFC3339, b.EndsAt)
		if err == nil && !now.Before(t) {
			return false
		}
	}
	return true
}

// ListPublic returns active banners in the optional slot, sorted for display.
func (b *BannerService) ListPublic(slot string) ([]models.Banner, error) {
	list, err := b.Store.ListBanners()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	want := normalizeBannerSlot(slot)
	var out []models.Banner
	for _, bn := range list {
		if want != "" && normalizeBannerSlot(bn.Slot) != want {
			continue
		}
		if !bannerVisibleAt(bn, now) {
			continue
		}
		out = append(out, bn)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].SortOrder != out[j].SortOrder {
			return out[i].SortOrder < out[j].SortOrder
		}
			return out[i].ID < out[j].ID
		})
	return out, nil
}

// AdminList returns all banners (no visibility filter).
func (b *BannerService) AdminList() ([]models.Banner, error) {
	list, err := b.Store.ListBanners()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].SortOrder != list[j].SortOrder {
			return list[i].SortOrder < list[j].SortOrder
		}
		return list[i].CreatedAt < list[j].CreatedAt
	})
	return list, nil
}

type BannerInput struct {
	Slot      string `json:"slot"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	ImageURL  string `json:"imageUrl"`
	LinkURL   string `json:"linkUrl"`
	SortOrder int    `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
	StartsAt  string `json:"startsAt"`
	EndsAt    string `json:"endsAt"`
}

func (b *BannerService) AdminCreate(in BannerInput) (*models.Banner, error) {
	slot := normalizeBannerSlot(in.Slot)
	if slot == "" {
		return nil, ErrValidation
	}
	now := time.Now().UTC().Format(time.RFC3339)
	bn := models.Banner{
		ID:        uuid.NewString(),
		Slot:      slot,
		Title:     strings.TrimSpace(in.Title),
		Body:      strings.TrimSpace(in.Body),
		ImageURL:  strings.TrimSpace(in.ImageURL),
		LinkURL:   strings.TrimSpace(in.LinkURL),
		SortOrder: in.SortOrder,
		IsActive:  in.IsActive,
		StartsAt:  strings.TrimSpace(in.StartsAt),
		EndsAt:    strings.TrimSpace(in.EndsAt),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := b.Store.UpsertBanner(bn); err != nil {
		return nil, err
	}
	return &bn, nil
}

func (b *BannerService) AdminUpdate(id string, in BannerInput) (*models.Banner, error) {
	cur, err := b.Store.FindBannerByID(id)
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, ErrNotFound
	}
	slot := normalizeBannerSlot(in.Slot)
	if slot == "" {
		slot = cur.Slot
	}
	now := time.Now().UTC().Format(time.RFC3339)
	cur.Slot = slot
	cur.Title = strings.TrimSpace(in.Title)
	cur.Body = strings.TrimSpace(in.Body)
	cur.ImageURL = strings.TrimSpace(in.ImageURL)
	cur.LinkURL = strings.TrimSpace(in.LinkURL)
	cur.SortOrder = in.SortOrder
	cur.IsActive = in.IsActive
	cur.StartsAt = strings.TrimSpace(in.StartsAt)
	cur.EndsAt = strings.TrimSpace(in.EndsAt)
	cur.UpdatedAt = now
	if err := b.Store.UpsertBanner(*cur); err != nil {
		return nil, err
	}
	return cur, nil
}

func (b *BannerService) AdminDelete(id string) error {
	cur, err := b.Store.FindBannerByID(id)
	if err != nil {
		return err
	}
	if cur == nil {
		return ErrNotFound
	}
	return b.Store.DeleteBanner(id)
}
