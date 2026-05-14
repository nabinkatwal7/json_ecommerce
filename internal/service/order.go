package service

import (
	"strings"
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	Store      *repository.Store
	Cart       *CartService
	Shipping   float64
	FreeShipAt float64
}

type CheckoutInput struct {
	ShippingAddress models.Address `json:"shippingAddress"`
	DiscountCode    string         `json:"discountCode"`
}

func (o *OrderService) Checkout(userID string, in CheckoutInput) (*models.Order, error) {
	cart, err := o.Cart.GetCart(userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, ErrValidation
	}
	if strings.TrimSpace(in.ShippingAddress.FullName) == "" ||
		strings.TrimSpace(in.ShippingAddress.AddressLine) == "" {
		return nil, ErrValidation
	}

	var items []models.OrderItem
	var subtotal float64
	for _, line := range cart.Items {
		p, err := o.Store.FindProductByID(line.ProductID)
		if err != nil {
			return nil, err
		}
		if p == nil || !p.IsActive {
			return nil, ErrValidation
		}
		var v *models.ProductVariant
		for i := range p.Variants {
			if p.Variants[i].ID == line.VariantID {
				v = &p.Variants[i]
				break
			}
		}
		if v == nil {
			return nil, ErrValidation
		}
		if v.Stock < line.Quantity {
			return nil, ErrInsufficientStock
		}
		// Price at checkout time from catalog (cart may be stale).
		price := v.Price
		items = append(items, models.OrderItem{
			ProductID: p.ID,
			VariantID: v.ID,
			Name:      p.Name,
			SKU:       v.SKU,
			Price:     price,
			Quantity:  line.Quantity,
		})
		subtotal += price * float64(line.Quantity)
	}

	now := time.Now().UTC()
	discountAmt := 0.0
	code := strings.TrimSpace(in.DiscountCode)
	if code != "" {
		d, err := o.Store.FindDiscountByCode(code)
		if err != nil {
			return nil, err
		}
		if d == nil {
			return nil, ErrNotFound
		}
		da, err := ApplyDiscount(d, subtotal, now)
		if err != nil {
			return nil, err
		}
		discountAmt = da
	}

	afterDiscount := subtotal - discountAmt
	if afterDiscount < 0 {
		afterDiscount = 0
	}
	shipping := o.Shipping
	if subtotal >= o.FreeShipAt {
		shipping = 0
	}
	total := afterDiscount + shipping

	order := models.Order{
		ID:              uuid.NewString(),
		UserID:          userID,
		Items:           items,
		ShippingAddress: in.ShippingAddress,
		Subtotal:        subtotal,
		Discount:        discountAmt,
		Shipping:        shipping,
		Total:           total,
		Status:          "created",
		PaymentStatus:   "pending",
		CreatedAt:       now.Format(time.RFC3339),
	}
	if err := o.Store.UpsertOrder(order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *OrderService) ListMyOrders(userID string) ([]models.Order, error) {
	return o.Store.OrdersByUser(userID)
}

func (o *OrderService) GetOrder(userID, orderID string) (*models.Order, error) {
	ord, err := o.Store.FindOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord == nil || ord.UserID != userID {
		return nil, ErrNotFound
	}
	return ord, nil
}

// Pay simulates a successful payment: captures funds, decrements inventory, marks order paid, clears cart.
func (o *OrderService) Pay(userID, orderID string) (*models.Order, *models.Payment, error) {
	ord, err := o.Store.FindOrderByID(orderID)
	if err != nil {
		return nil, nil, err
	}
	if ord == nil || ord.UserID != userID {
		return nil, nil, ErrNotFound
	}
	if ord.PaymentStatus != "pending" || ord.Status != "created" {
		return nil, nil, ErrBadState
	}

	products, err := o.Store.ListProducts()
	if err != nil {
		return nil, nil, err
	}
	byID := make(map[string]int)
	for i := range products {
		byID[products[i].ID] = i
	}

	for _, line := range ord.Items {
		idx, ok := byID[line.ProductID]
		if !ok {
			return nil, nil, ErrValidation
		}
		p := &products[idx]
		if !p.IsActive {
			return nil, nil, ErrInactive
		}
		vidx := -1
		for i := range p.Variants {
			if p.Variants[i].ID == line.VariantID {
				vidx = i
				break
			}
		}
		if vidx < 0 {
			return nil, nil, ErrValidation
		}
		if p.Variants[vidx].Stock < line.Quantity {
			return nil, nil, ErrInsufficientStock
		}
		p.Variants[vidx].Stock -= line.Quantity
		p.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if err := o.Store.SaveProducts(products); err != nil {
		return nil, nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	pay := models.Payment{
		ID:                uuid.NewString(),
		OrderID:           ord.ID,
		Provider:          "stub",
		ProviderReference: uuid.NewString(),
		Amount:            ord.Total,
		Status:            "paid",
		CreatedAt:         now,
	}
	if err := o.Store.UpsertPayment(pay); err != nil {
		return nil, nil, err
	}

	ord.Status = "paid"
	ord.PaymentStatus = "paid"
	if err := o.Store.UpsertOrder(*ord); err != nil {
		return nil, nil, err
	}
	if err := o.Cart.ClearCartByUser(userID); err != nil {
		return nil, nil, err
	}
	return ord, &pay, nil
}
