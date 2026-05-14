package invoice

import (
	"bytes"
	"fmt"
	"strings"

	"go-ecommerce-json/internal/models"

	"github.com/go-pdf/fpdf"
)

// BuildOrderPDF renders a simple commercial invoice (no payment card data).
func BuildOrderPDF(order *models.Order) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, "Invoice")
	pdf.Ln(12)

	pdf.SetFont("Helvetica", "", 11)
	inv := order.InvoiceNumber
	if inv == "" {
		inv = order.ID
	}
	pdf.Cell(0, 6, fmt.Sprintf("Invoice #: %s", inv))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Order #: %s", order.ID))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Date: %s", order.CreatedAt))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Status: %s  Payment: %s", order.Status, order.PaymentStatus))
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 8, "Ship to")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 10)
	addr := order.ShippingAddress
	lines := []string{
		addr.FullName,
		addr.AddressLine,
		fmt.Sprintf("%s, %s %s", addr.City, addr.State, addr.PostalCode),
		addr.Country,
		"Phone: " + addr.Phone,
	}
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		pdf.Cell(0, 5, ln)
		pdf.Ln(5)
	}
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(80, 7, "Item", "1", 0, "L", false, 0, "")
	pdf.CellFormat(25, 7, "SKU", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 7, "Qty", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, "Price", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, "Line", "1", 1, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	for _, it := range order.Items {
		lineTotal := it.Price * float64(it.Quantity)
		pdf.CellFormat(80, 7, truncate(it.Name, 46), "LR", 0, "L", false, 0, "")
		pdf.CellFormat(25, 7, truncate(it.SKU, 12), "LR", 0, "L", false, 0, "")
		pdf.CellFormat(20, 7, fmt.Sprintf("%d", it.Quantity), "LR", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", it.Price), "LR", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", lineTotal), "LR", 1, "R", false, 0, "")
	}
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(155, 7, "Subtotal", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", order.Subtotal), "0", 1, "R", false, 0, "")
	pdf.CellFormat(155, 7, "Discount", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", order.Discount), "0", 1, "R", false, 0, "")
	pdf.CellFormat(155, 7, "Shipping", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", order.Shipping), "0", 1, "R", false, 0, "")
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(155, 8, "Total", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", order.Total), "0", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
