package service

import (
	"fmt"
	"strings"

	"go-ecommerce-json/internal/repository"
)

type LowStockLine struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	VariantID   string `json:"variantId"`
	SKU         string `json:"sku"`
	Stock       int    `json:"stock"`
}

func LowStockReport(store *repository.Store, threshold int) ([]LowStockLine, error) {
	if threshold < 0 {
		threshold = 0
	}
	products, err := store.ListProducts()
	if err != nil {
		return nil, err
	}
	var out []LowStockLine
	for _, p := range products {
		if !p.IsActive {
			continue
		}
		for _, v := range p.Variants {
			if v.Stock <= threshold {
				out = append(out, LowStockLine{
					ProductID:   p.ID,
					ProductName: p.Name,
					VariantID:   v.ID,
					SKU:         v.SKU,
					Stock:       v.Stock,
				})
			}
		}
	}
	return out, nil
}

func FormatLowStockEmail(lines []LowStockLine) string {
	if len(lines) == 0 {
		return "No variants at or below the threshold."
	}
	var b strings.Builder
	b.WriteString("Low stock alert:\n\n")
	for _, ln := range lines {
		b.WriteString(fmt.Sprintf("- %s / SKU %s / stock=%d (product %s)\n", ln.ProductName, ln.SKU, ln.Stock, ln.ProductID))
	}
	return b.String()
}
