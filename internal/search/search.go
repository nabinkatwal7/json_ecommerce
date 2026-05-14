package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"
)

// Client powers catalog search (Algolia when configured, otherwise local fuzzy ranking).
type Client struct {
	Store         *repository.Store
	AlgoliaAppID  string
	AlgoliaAPIKey string
	AlgoliaIndex  string
	HTTP          *http.Client
}

type Hit struct {
	ObjectID   string   `json:"objectID"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	CategoryID string   `json:"categoryId"`
	Tags       []string `json:"tags"`
	MinPrice   float64  `json:"minPrice"`
	Score      int      `json:"score,omitempty"`
}

type algoliaDoc struct {
	ObjectID   string   `json:"objectID"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	CategoryID string   `json:"categoryId"`
	Tags       []string `json:"tags"`
	MinPrice   float64  `json:"minPrice"`
	Active     bool     `json:"active"`
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return &http.Client{Timeout: 8 * time.Second}
}

func (c *Client) algoliaEnabled() bool {
	return strings.TrimSpace(c.AlgoliaAppID) != "" &&
		strings.TrimSpace(c.AlgoliaAPIKey) != "" &&
		strings.TrimSpace(c.AlgoliaIndex) != ""
}

// Search runs Algolia or local fuzzy search.
func (c *Client) Search(q, categoryID string, limit int) ([]Hit, error) {
	q = strings.TrimSpace(q)
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if c.algoliaEnabled() {
		return c.searchAlgolia(q, categoryID, limit)
	}
	return c.searchLocal(q, categoryID, limit)
}

func (c *Client) searchAlgolia(q, categoryID string, limit int) ([]Hit, error) {
	host := fmt.Sprintf("https://%s-dsn.algolia.net", strings.TrimSpace(c.AlgoliaAppID))
	u := host + "/1/indexes/" + url.PathEscape(c.AlgoliaIndex) + "/query"
	params := "hitsPerPage=" + fmt.Sprint(limit) + "&query=" + url.QueryEscape(q)
	if categoryID != "" {
		params += "&facetFilters=" + url.QueryEscape("categoryId:"+categoryID)
	}
	payload := map[string]string{"params": params}
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Algolia-Application-Id", c.AlgoliaAppID)
	req.Header.Set("X-Algolia-API-Key", c.AlgoliaAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("algolia: %s: %s", resp.Status, string(b))
	}
	var out struct {
		Hits []Hit `json:"hits"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out.Hits, nil
}

func (c *Client) searchLocal(q, categoryID string, limit int) ([]Hit, error) {
	if q == "" {
		return nil, nil
	}
	all, err := c.Store.ListProducts()
	if err != nil {
		return nil, err
	}
	var haystacks []string
	var products []*models.Product
	for i := range all {
		p := &all[i]
		if !p.IsActive {
			continue
		}
		if categoryID != "" && p.CategoryID != categoryID {
			continue
		}
		h := strings.ToLower(p.Name + " " + p.Slug + " " + strings.Join(p.Tags, " "))
		haystacks = append(haystacks, h)
		products = append(products, p)
	}
	ranks := fuzzy.RankFindFold(q, haystacks)
	var hits []Hit
	for _, r := range ranks {
		if len(hits) >= limit {
			break
		}
		if r.OriginalIndex < 0 || r.OriginalIndex >= len(products) {
			continue
		}
		p := products[r.OriginalIndex]
		hits = append(hits, Hit{
			ObjectID:   p.ID,
			Name:       p.Name,
			Slug:       p.Slug,
			CategoryID: p.CategoryID,
			Tags:       p.Tags,
			MinPrice:   minVariantPrice(p),
			Score:      r.Distance,
		})
	}
	return hits, nil
}

func minVariantPrice(p *models.Product) float64 {
	if len(p.Variants) == 0 {
		return 0
	}
	m := p.Variants[0].Price
	for _, v := range p.Variants[1:] {
		if v.Price < m {
			m = v.Price
		}
	}
	return m
}

// ReindexAlgolia pushes active products to Algolia using the batch API (chunked).
func (c *Client) ReindexAlgolia() (int, error) {
	if !c.algoliaEnabled() {
		return 0, fmt.Errorf("algolia not configured")
	}
	all, err := c.Store.ListProducts()
	if err != nil {
		return 0, err
	}
	var docs []algoliaDoc
	for _, p := range all {
		if !p.IsActive {
			continue
		}
		docs = append(docs, algoliaDoc{
			ObjectID:   p.ID,
			Name:       p.Name,
			Slug:       p.Slug,
			CategoryID: p.CategoryID,
			Tags:       p.Tags,
			MinPrice:   minVariantPrice(&p),
			Active:     true,
		})
	}
	host := fmt.Sprintf("https://%s-dsn.algolia.net", strings.TrimSpace(c.AlgoliaAppID))
	u := host + "/1/indexes/" + url.PathEscape(c.AlgoliaIndex) + "/batch"
	const chunk = 500
	total := 0
	for i := 0; i < len(docs); i += chunk {
		j := i + chunk
		if j > len(docs) {
			j = len(docs)
		}
		part := docs[i:j]
		var reqs []map[string]any
		for _, d := range part {
			o := d
			reqs = append(reqs, map[string]any{
				"action": "partialUpdateObject",
				"body":   o,
			})
		}
		raw, _ := json.Marshal(map[string]any{"requests": reqs})
		req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(raw))
		if err != nil {
			return total, err
		}
		req.Header.Set("X-Algolia-Application-Id", c.AlgoliaAppID)
		req.Header.Set("X-Algolia-API-Key", c.AlgoliaAPIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.httpClient().Do(req)
		if err != nil {
			return total, err
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return total, fmt.Errorf("algolia batch: %s: %s", resp.Status, string(b))
		}
		total += len(part)
	}
	return total, nil
}
