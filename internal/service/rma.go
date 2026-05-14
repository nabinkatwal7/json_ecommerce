package service

import (
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type RMAService struct {
	Store *repository.Store
}

type RMACreateItem struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
	Quantity  int    `json:"quantity"`
}

type RMACreateInput struct {
	OrderID string `json:"orderId"`
	Reason  string `json:"reason"`
	Items   []RMACreateItem `json:"items"`
}

func orderLineQty(ord *models.Order, productID, variantID string) int {
	for _, it := range ord.Items {
		if it.ProductID == productID && it.VariantID == variantID {
			return it.Quantity
		}
	}
	return 0
}

func (r *RMAService) Create(userID string, in RMACreateInput) (*models.RMA, error) {
	orderID := strings.TrimSpace(in.OrderID)
	if strings.TrimSpace(in.Reason) == "" || len(in.Items) == 0 || orderID == "" {
		return nil, ErrValidation
	}
	ord, err := r.Store.FindOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord == nil || ord.UserID != userID {
		return nil, ErrNotFound
	}
	if ord.PaymentStatus != "paid" {
		return nil, ErrBadState
	}
	if ord.Status != "fulfilled" && ord.Status != "shipped" {
		return nil, ErrBadState
	}
	var ritems []models.RMAItem
	for _, req := range in.Items {
		if req.Quantity <= 0 {
			return nil, ErrValidation
		}
		max := orderLineQty(ord, req.ProductID, req.VariantID)
		if max == 0 || req.Quantity > max {
			return nil, ErrValidation
		}
		var name, sku string
		var price float64
		for _, it := range ord.Items {
			if it.ProductID == req.ProductID && it.VariantID == req.VariantID {
				name, sku, price = it.Name, it.SKU, it.Price
				break
			}
		}
		ritems = append(ritems, models.RMAItem{
			ProductID: req.ProductID,
			VariantID: req.VariantID,
			SKU:       sku,
			Name:      name,
			Quantity:  req.Quantity,
			Price:     price,
		})
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	rm := models.RMA{
		ID:        uuid.NewString(),
		UserID:    userID,
		OrderID:   orderID,
		Items:     ritems,
		Reason:    strings.TrimSpace(in.Reason),
		Status:    "requested",
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	if err := r.Store.UpsertRMA(rm); err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *RMAService) ListMine(userID string) ([]models.RMA, error) {
	return r.Store.RMAsByUser(userID)
}

func (r *RMAService) GetMine(userID, id string) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil {
		return nil, err
	}
	if rm == nil || rm.UserID != userID {
		return nil, ErrNotFound
	}
	return rm, nil
}

func (r *RMAService) ListAdmin() ([]models.RMA, error) {
	return r.Store.ListRMAs()
}

func (r *RMAService) GetAdmin(id string) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil {
		return nil, err
	}
	if rm == nil {
		return nil, ErrNotFound
	}
	return rm, nil
}

func (r *RMAService) setStatus(id, status, note string, refund float64) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil {
		return nil, err
	}
	if rm == nil {
		return nil, ErrNotFound
	}
	rm.Status = status
	rm.AdminNote = note
	rm.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if refund > 0 {
		rm.RefundAmount = refund
	}
	if err := r.Store.UpsertRMA(*rm); err != nil {
		return nil, err
	}
	return rm, nil
}

func (r *RMAService) AdminApprove(id, note string) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil || rm == nil {
		return nil, ErrNotFound
	}
	if rm.Status != "requested" {
		return nil, ErrBadState
	}
	return r.setStatus(id, "approved", note, 0)
}

func (r *RMAService) AdminReject(id, note string) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil || rm == nil {
		return nil, ErrNotFound
	}
	if rm.Status != "requested" && rm.Status != "approved" {
		return nil, ErrBadState
	}
	return r.setStatus(id, "rejected", note, 0)
}

func (r *RMAService) AdminMarkReceived(id, note string) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil || rm == nil {
		return nil, ErrNotFound
	}
	if rm.Status != "approved" {
		return nil, ErrBadState
	}
	return r.setStatus(id, "received", note, 0)
}

// AdminRefund marks refunded and restores inventory for returned units (payment gateway refund is separate).
func (r *RMAService) AdminRefund(id, note string) (*models.RMA, error) {
	rm, err := r.Store.FindRMAByID(id)
	if err != nil || rm == nil {
		return nil, ErrNotFound
	}
	if rm.Status != "received" {
		return nil, ErrBadState
	}
	var oitems []models.OrderItem
	for _, it := range rm.Items {
		oitems = append(oitems, models.OrderItem{
			ProductID: it.ProductID,
			VariantID: it.VariantID,
			Quantity:  it.Quantity,
		})
	}
	if err := adjustVariantStock(r.Store, oitems, +1); err != nil {
		return nil, err
	}
	var refund float64
	for _, it := range rm.Items {
		refund += it.Price * float64(it.Quantity)
	}
	return r.setStatus(id, "refunded", note, refund)
}
