package service

import (
	"sort"
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"
)

type DiscoveryService struct {
	Store *repository.Store
}

type FeaturedCategoryRow struct {
	models.Category
	ProductCount int `json:"productCount"`
}

type FeaturedHome struct {
	NewArrivals        []models.Product        `json:"newArrivals"`
	BestSellers        []models.Product        `json:"bestSellers"`
	FeaturedCategories []FeaturedCategoryRow   `json:"featuredCategories"`
}

func parseRFC3339OrZero(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// FeaturedHome returns a single payload for storefront home pages (new, popular, categories with counts).
func (d *DiscoveryService) FeaturedHome() (*FeaturedHome, error) {
	products, err := d.Store.ListProducts()
	if err != nil {
		return nil, err
	}
	var active []models.Product
	for _, p := range products {
		if p.IsActive {
			active = append(active, p)
		}
	}
	newArr := append([]models.Product(nil), active...)
	sort.Slice(newArr, func(i, j int) bool {
		ti := parseRFC3339OrZero(newArr[i].CreatedAt)
		tj := parseRFC3339OrZero(newArr[j].CreatedAt)
		return ti.After(tj)
	})
	if len(newArr) > 8 {
		newArr = newArr[:8]
	}

	orders, err := d.Store.ListOrders()
	if err != nil {
		return nil, err
	}
	qtyByProduct := map[string]int{}
	for _, o := range orders {
		if strings.EqualFold(o.Status, "cancelled") {
			continue
		}
		if !strings.EqualFold(o.PaymentStatus, "paid") {
			continue
		}
		for _, it := range o.Items {
			qtyByProduct[it.ProductID] += it.Quantity
		}
	}
	type scored struct {
		id  string
		qty int
	}
	var bests []scored
	for id, q := range qtyByProduct {
		bests = append(bests, scored{id: id, qty: q})
	}
	sort.Slice(bests, func(i, j int) bool {
		if bests[i].qty != bests[j].qty {
			return bests[i].qty > bests[j].qty
		}
		return bests[i].id < bests[j].id
	})
	var bestProducts []models.Product
	byID := map[string]models.Product{}
	for _, p := range active {
		byID[p.ID] = p
	}
	for _, b := range bests {
		if p, ok := byID[b.id]; ok && len(bestProducts) < 8 {
			bestProducts = append(bestProducts, p)
		}
	}
	if len(bestProducts) < 4 && len(active) > 0 {
		for _, p := range active {
			if qtyByProduct[p.ID] == 0 && len(bestProducts) < 8 {
				bestProducts = append(bestProducts, p)
			}
			if len(bestProducts) >= 8 {
				break
			}
		}
	}

	cats, err := d.Store.ListCategories()
	if err != nil {
		return nil, err
	}
	var feat []FeaturedCategoryRow
	for _, c := range cats {
		if !c.IsActive {
			continue
		}
		n := 0
		for _, p := range active {
			if p.CategoryID == c.ID {
				n++
			}
		}
		feat = append(feat, FeaturedCategoryRow{Category: c, ProductCount: n})
	}
	sort.Slice(feat, func(i, j int) bool {
		if feat[i].ProductCount != feat[j].ProductCount {
			return feat[i].ProductCount > feat[j].ProductCount
		}
		return feat[i].Name < feat[j].Name
	})
	if len(feat) > 6 {
		feat = feat[:6]
	}

	return &FeaturedHome{
		NewArrivals:        newArr,
		BestSellers:        bestProducts,
		FeaturedCategories: feat,
	}, nil
}

// RelatedProducts returns similar active products (category + tag overlap), excluding the seed product.
func (d *DiscoveryService) RelatedProducts(productID string, limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 8
	}
	if limit > 8 {
		limit = 8
	}
	seed, err := d.Store.FindProductByID(productID)
	if err != nil {
		return nil, err
	}
	if seed == nil || !seed.IsActive {
		return nil, ErrNotFound
	}
	tagSet := map[string]struct{}{}
	for _, id := range seed.TagIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			tagSet[id] = struct{}{}
		}
	}
	products, err := d.Store.ListProducts()
	if err != nil {
		return nil, err
	}
	type cand struct {
		p     models.Product
		score int
	}
	var cands []cand
	for _, p := range products {
		if !p.IsActive || p.ID == seed.ID {
			continue
		}
		sc := 0
		if p.CategoryID != "" && p.CategoryID == seed.CategoryID {
			sc += 2
		}
		for _, tid := range p.TagIDs {
			tid = strings.TrimSpace(tid)
			if tid == "" {
				continue
			}
			if _, ok := tagSet[tid]; ok {
				sc++
			}
		}
		if sc == 0 {
			continue
		}
		cands = append(cands, cand{p: p, score: sc})
	}
	if len(cands) == 0 {
		for _, p := range products {
			if !p.IsActive || p.ID == seed.ID {
				continue
			}
			if p.CategoryID != "" && p.CategoryID == seed.CategoryID {
				cands = append(cands, cand{p: p, score: 1})
			}
		}
	}
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].score != cands[j].score {
			return cands[i].score > cands[j].score
		}
		ti := parseRFC3339OrZero(cands[i].p.CreatedAt)
		tj := parseRFC3339OrZero(cands[j].p.CreatedAt)
		return ti.After(tj)
	})
	var out []models.Product
	for _, c := range cands {
		out = append(out, c.p)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

// CatalogCounts returns how many active products and active categories exist.
func (d *DiscoveryService) CatalogCounts() (productCount int, categoryCount int, err error) {
	products, err := d.Store.ListProducts()
	if err != nil {
		return 0, 0, err
	}
	for _, p := range products {
		if p.IsActive {
			productCount++
		}
	}
	cats, err := d.Store.ListCategories()
	if err != nil {
		return 0, 0, err
	}
	for _, c := range cats {
		if c.IsActive {
			categoryCount++
		}
	}
	return productCount, categoryCount, nil
}

// SaleProducts lists active products that carry the "sale" tag (by slug), newest first.
func (d *DiscoveryService) SaleProducts(limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 8
	}
	if limit > 24 {
		limit = 24
	}
	tags, err := d.Store.ListTags()
	if err != nil {
		return nil, err
	}
	var saleTagID string
	for _, t := range tags {
		if strings.EqualFold(strings.TrimSpace(t.Slug), "sale") {
			saleTagID = t.ID
			break
		}
	}
	if saleTagID == "" {
		return []models.Product{}, nil
	}
	products, err := d.Store.ListProducts()
	if err != nil {
		return nil, err
	}
	var matches []models.Product
	for _, p := range products {
		if !p.IsActive {
			continue
		}
		for _, tid := range p.TagIDs {
			if strings.TrimSpace(tid) == saleTagID {
				matches = append(matches, p)
				break
			}
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		ti := parseRFC3339OrZero(matches[i].UpdatedAt)
		tj := parseRFC3339OrZero(matches[j].UpdatedAt)
		if ti.Equal(tj) {
			return matches[i].Name < matches[j].Name
		}
		return ti.After(tj)
	})
	if len(matches) > limit {
		matches = matches[:limit]
	}
	return matches, nil
}

// TagIDBySlug returns the tag id for a slug, or empty string if unknown.
func (d *DiscoveryService) TagIDBySlug(want string) string {
	want = strings.TrimSpace(strings.ToLower(want))
	if want == "" {
		return ""
	}
	tags, err := d.Store.ListTags()
	if err != nil {
		return ""
	}
	for _, t := range tags {
		if strings.EqualFold(strings.TrimSpace(t.Slug), want) {
			return t.ID
		}
	}
	return ""
}

// SearchSuggestions returns distinct product names matching a short query (substring, case-insensitive).
func (d *DiscoveryService) SearchSuggestions(q string, limit int) ([]string, error) {
	q = strings.TrimSpace(strings.ToLower(q))
	if len(q) < 2 {
		return []string{}, nil
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}
	products, err := d.Store.ListProducts()
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	var out []string
	for _, p := range products {
		if !p.IsActive {
			continue
		}
		name := strings.TrimSpace(p.Name)
		if name == "" {
			continue
		}
		if !strings.Contains(strings.ToLower(name), q) {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
		if len(out) >= limit {
			break
		}
	}
	sort.Strings(out)
	return out, nil
}
