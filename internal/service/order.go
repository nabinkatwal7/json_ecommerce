package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-ecommerce-json/internal/invoice"
	"go-ecommerce-json/internal/mail"
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/payment"
	"go-ecommerce-json/internal/repository"
	"go-ecommerce-json/internal/shipping"

	"github.com/google/uuid"
)

type OrderService struct {
	Store               *repository.Store
	Cart                *CartService
	Shipping            float64
	FreeShipAt          float64
	DefaultItemWeightKg float64
	Mail                *mail.Sender
	StripeSecret        string
	StripeCurrency      string
	DevPaymentStub      bool
	AppPublicURL        string
	LowStockThreshold   int
	AdminAlertEmail     string
}

type CheckoutInput struct {
	ShippingAddress models.Address `json:"shippingAddress"`
	DiscountCode    string         `json:"discountCode"`
	// ShippingCarrier: empty or "flat" uses flat rate + free-ship threshold; otherwise fedex_ground | ups_ground | dhl_express (stub quotes).
	ShippingCarrier string `json:"shippingCarrier"`
}

type PayInput struct {
	// StripePaymentIntentID is created client-side with Stripe.js/Elements (card never hits this API).
	StripePaymentIntentID string `json:"stripePaymentIntentId"`
	Stub                  bool   `json:"stub"`
}

// CouponValidateResult is a checkout-friendly coupon check (no order created).
type CouponValidateResult struct {
	Valid          bool    `json:"valid"`
	Code           string  `json:"code,omitempty"`
	Message        string  `json:"message"`
	DiscountType   string  `json:"discountType,omitempty"`
	DiscountAmount float64 `json:"discountAmount"`
	Subtotal       float64 `json:"subtotal"`
}

func (o *OrderService) buildOrderItemsFromCart(userID string) ([]models.OrderItem, float64, float64, error) {
	cart, err := o.Cart.GetCart(userID)
	if err != nil {
		return nil, 0, 0, err
	}
	var items []models.OrderItem
	var subtotal float64
	var totalWeight float64
	for _, line := range cart.Items {
		p, err := o.Store.FindProductByID(line.ProductID)
		if err != nil {
			return nil, 0, 0, err
		}
		if p == nil || !p.IsActive {
			return nil, 0, 0, ErrValidation
		}
		var v *models.ProductVariant
		for i := range p.Variants {
			if p.Variants[i].ID == line.VariantID {
				v = &p.Variants[i]
				break
			}
		}
		if v == nil {
			return nil, 0, 0, ErrValidation
		}
		if v.Stock < line.Quantity {
			return nil, 0, 0, ErrInsufficientStock
		}
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
		w := v.WeightKg
		if w <= 0 {
			w = o.DefaultItemWeightKg
			if w <= 0 {
				w = 0.5
			}
		}
		totalWeight += w * float64(line.Quantity)
	}
	return items, subtotal, totalWeight, nil
}

// ValidateCouponForCart checks a promo code against the current cart (same rules as checkout).
func (o *OrderService) ValidateCouponForCart(userID, code string) (*CouponValidateResult, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return &CouponValidateResult{Valid: false, Message: "code required"}, nil
	}
	items, subtotal, _, err := o.buildOrderItemsFromCart(userID)
	if err != nil {
		if errors.Is(err, ErrInsufficientStock) {
			return &CouponValidateResult{Valid: false, Code: code, Message: "insufficient stock for cart", Subtotal: subtotal}, nil
		}
		if errors.Is(err, ErrValidation) {
			return &CouponValidateResult{Valid: false, Code: code, Message: "cart has invalid or inactive items", Subtotal: subtotal}, nil
		}
		return nil, err
	}
	if len(items) == 0 {
		return &CouponValidateResult{Valid: false, Code: code, Message: "cart is empty", Subtotal: 0}, nil
	}
	d, err := o.Store.FindDiscountByCode(code)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return &CouponValidateResult{Valid: false, Code: code, Message: "unknown code", Subtotal: subtotal}, nil
	}
	amt, err := ComputeDiscountAmount(o.Store, d, subtotal, items, time.Now().UTC())
	if err != nil {
		msg := "not applicable"
		switch {
		case errors.Is(err, ErrInactive):
			msg = "inactive or expired"
		case errors.Is(err, ErrValidation):
			msg = "minimum not met or invalid for cart"
		default:
			msg = "not applicable"
		}
		return &CouponValidateResult{Valid: false, Code: code, Message: msg, DiscountType: d.Type, Subtotal: subtotal}, nil
	}
	return &CouponValidateResult{
		Valid:          true,
		Code:           code,
		Message:        "ok",
		DiscountType:   d.Type,
		DiscountAmount: amt,
		Subtotal:       subtotal,
	}, nil
}

func (o *OrderService) Checkout(userID string, in CheckoutInput) (*models.Order, error) {
	if strings.TrimSpace(in.ShippingAddress.FullName) == "" ||
		strings.TrimSpace(in.ShippingAddress.AddressLine) == "" {
		return nil, ErrValidation
	}

	items, subtotal, totalWeight, err := o.buildOrderItemsFromCart(userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrValidation
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
		da, err := ComputeDiscountAmount(o.Store, d, subtotal, items, now)
		if err != nil {
			return nil, err
		}
		discountAmt = da
	}

	afterDiscount := subtotal - discountAmt
	if afterDiscount < 0 {
		afterDiscount = 0
	}
	shipLabel := "flat"
	shippingAmt := o.Shipping
	carrierIn := strings.TrimSpace(strings.ToLower(in.ShippingCarrier))
	if carrierIn == "" || carrierIn == "flat" {
		if subtotal >= o.FreeShipAt {
			shippingAmt = 0
		}
	} else {
		q, ok := shipping.FindQuote(in.ShippingCarrier, in.ShippingAddress.Country, totalWeight)
		if !ok {
			return nil, ErrValidation
		}
		shippingAmt = q.Amount
		shipLabel = q.Code
		if subtotal >= o.FreeShipAt {
			shippingAmt = 0
		}
	}
	total := afterDiscount + shippingAmt

	ts := now.Format(time.RFC3339)
	order := models.Order{
		ID:              uuid.NewString(),
		UserID:          userID,
		Items:           items,
		ShippingAddress: in.ShippingAddress,
		Subtotal:        subtotal,
		Discount:        discountAmt,
		Shipping:        shippingAmt,
		ShippingCarrier: shipLabel,
		Total:           total,
		Status:          "created",
		PaymentStatus:   "pending",
		CreatedAt:       ts,
		UpdatedAt:       ts,
	}

	if err := adjustVariantStock(o.Store, items, -1); err != nil {
		return nil, err
	}
	if err := o.Store.UpsertOrder(order); err != nil {
		_ = adjustVariantStock(o.Store, items, +1)
		return nil, err
	}
	_ = o.Cart.ClearCartByUser(userID)

	o.maybeLowStockAlert()
	return &order, nil
}

func (o *OrderService) maybeLowStockAlert() {
	th := o.LowStockThreshold
	if th <= 0 {
		th = 5
	}
	if o.Mail == nil || !o.Mail.Configured() || strings.TrimSpace(o.AdminAlertEmail) == "" {
		return
	}
	lines, err := LowStockReport(o.Store, th)
	if err != nil || len(lines) == 0 {
		return
	}
	subj := fmt.Sprintf("Low stock alert (%d)", len(lines))
	_ = o.Mail.SendPlain(o.AdminAlertEmail, subj, FormatLowStockEmail(lines))
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

func (o *OrderService) GetOrderAdmin(orderID string) (*models.Order, error) {
	ord, err := o.Store.FindOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord == nil {
		return nil, ErrNotFound
	}
	return ord, nil
}

func (o *OrderService) ListOrdersAdmin() ([]models.Order, error) {
	return o.Store.ListOrders()
}

func (o *OrderService) touchOrder(ord *models.Order) {
	ord.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func (o *OrderService) ensureInvoiceNumber(ord *models.Order) error {
	if ord.InvoiceNumber != "" {
		return nil
	}
	ord.InvoiceNumber = fmt.Sprintf("INV-%s-%s", time.Now().UTC().Format("20060102"), ord.ID[:8])
	o.touchOrder(ord)
	return o.Store.UpsertOrder(*ord)
}

// InvoicePDF returns a PDF for the order owner. Assigns an invoice number on first generation.
func (o *OrderService) InvoicePDF(userID, orderID string) ([]byte, error) {
	ord, err := o.GetOrder(userID, orderID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(ord.Status, "cancelled") {
		return nil, ErrBadState
	}
	if err := o.ensureInvoiceNumber(ord); err != nil {
		return nil, err
	}
	ord, err = o.GetOrder(userID, orderID)
	if err != nil {
		return nil, err
	}
	return invoice.BuildOrderPDF(ord)
}

// CancelByCustomer cancels a pending, unpaid order and restores inventory.
func (o *OrderService) CancelByCustomer(userID, orderID string) (*models.Order, error) {
	ord, err := o.GetOrder(userID, orderID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(ord.Status, "cancelled") {
		return nil, ErrBadState
	}
	if ord.Status != "created" || ord.PaymentStatus != "pending" {
		return nil, ErrBadState
	}
	if err := adjustVariantStock(o.Store, ord.Items, +1); err != nil {
		return nil, err
	}
	ord.Status = "cancelled"
	ord.PaymentStatus = "failed"
	now := time.Now().UTC().Format(time.RFC3339)
	ord.CancelledAt = now
	o.touchOrder(ord)
	if err := o.Store.UpsertOrder(*ord); err != nil {
		_ = adjustVariantStock(o.Store, ord.Items, -1)
		return nil, err
	}
	return ord, nil
}

// AdminCancel cancels an order (except already cancelled) and restores inventory once.
func (o *OrderService) AdminCancel(orderID string) (*models.Order, error) {
	ord, err := o.GetOrderAdmin(orderID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(ord.Status, "cancelled") {
		return nil, ErrBadState
	}
	// Restore stock for any order that had checkout deduction (all non-cancelled prior states).
	if err := adjustVariantStock(o.Store, ord.Items, +1); err != nil {
		return nil, err
	}
	ord.Status = "cancelled"
	if ord.PaymentStatus == "paid" {
		ord.PaymentStatus = "refunded"
	} else {
		ord.PaymentStatus = "failed"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	ord.CancelledAt = now
	o.touchOrder(ord)
	if err := o.Store.UpsertOrder(*ord); err != nil {
		_ = adjustVariantStock(o.Store, ord.Items, -1)
		return nil, err
	}
	return ord, nil
}

func (o *OrderService) AdminFulfill(orderID string) (*models.Order, error) {
	ord, err := o.GetOrderAdmin(orderID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(ord.Status, "cancelled") {
		return nil, ErrBadState
	}
	if ord.PaymentStatus != "paid" || ord.Status != "paid" {
		return nil, ErrBadState
	}
	ord.Status = "fulfilled"
	now := time.Now().UTC().Format(time.RFC3339)
	ord.FulfilledAt = now
	o.touchOrder(ord)
	if err := o.Store.UpsertOrder(*ord); err != nil {
		return nil, err
	}
	o.sendOrderEmail(ord, "Your order has been fulfilled", "We are preparing your shipment.\n\nOrder: "+ord.ID)
	return ord, nil
}

func (o *OrderService) AdminShip(orderID string) (*models.Order, error) {
	ord, err := o.GetOrderAdmin(orderID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(ord.Status, "cancelled") {
		return nil, ErrBadState
	}
	if ord.Status != "fulfilled" {
		return nil, ErrBadState
	}
	ord.Status = "shipped"
	now := time.Now().UTC().Format(time.RFC3339)
	ord.ShippedAt = now
	o.touchOrder(ord)
	if err := o.Store.UpsertOrder(*ord); err != nil {
		return nil, err
	}
	o.sendOrderEmail(ord, "Your order has shipped", "Your order is on the way.\n\nOrder: "+ord.ID)
	return ord, nil
}

func (o *OrderService) sendOrderEmail(ord *models.Order, subject, body string) {
	if o.Mail == nil || !o.Mail.Configured() {
		return
	}
	u, err := o.Store.FindUserByID(ord.UserID)
	if err != nil || u == nil {
		return
	}
	_ = o.Mail.SendPlain(u.Email, subject, body)
}

// Pay confirms payment (Stripe PaymentIntent or dev stub). Inventory was already reduced at checkout.
func (o *OrderService) Pay(userID, orderID string, in PayInput) (*models.Order, *models.Payment, error) {
	ord, err := o.GetOrder(userID, orderID)
	if err != nil {
		return nil, nil, err
	}
	if ord.PaymentStatus != "pending" || ord.Status != "created" {
		return nil, nil, ErrBadState
	}

	cur := strings.ToLower(strings.TrimSpace(o.StripeCurrency))
	if cur == "" {
		cur = "usd"
	}

	var providerRef string
	var provider string

	switch {
	case strings.TrimSpace(o.StripeSecret) != "":
		if strings.TrimSpace(in.StripePaymentIntentID) == "" {
			return nil, nil, ErrValidation
		}
		if err := payment.VerifySucceededPaymentIntent(o.StripeSecret, in.StripePaymentIntentID, ord.Total, cur); err != nil {
			return nil, nil, ErrPayment
		}
		provider = "stripe"
		providerRef = in.StripePaymentIntentID
	case o.DevPaymentStub && in.Stub:
		provider = "stub"
		providerRef = uuid.NewString()
	default:
		return nil, nil, ErrValidation
	}

	now := time.Now().UTC().Format(time.RFC3339)
	pay := models.Payment{
		ID:                uuid.NewString(),
		OrderID:           ord.ID,
		Provider:          provider,
		ProviderReference: providerRef,
		Amount:            ord.Total,
		Status:            "paid",
		CreatedAt:         now,
	}
	if err := o.Store.UpsertPayment(pay); err != nil {
		return nil, nil, err
	}
	ord.Status = "paid"
	ord.PaymentStatus = "paid"
	ord.PaidAt = now
	o.touchOrder(ord)
	if err := o.Store.UpsertOrder(*ord); err != nil {
		return nil, nil, err
	}
	body := fmt.Sprintf("Payment received for order %s.\nTotal: %.2f\nThank you.\n", ord.ID, ord.Total)
	o.sendOrderEmail(ord, "Order paid — thank you", body)
	return ord, &pay, nil
}

// CreateStripePaymentIntent creates a server-side PaymentIntent so the client never chooses the charge amount.
func (o *OrderService) CreateStripePaymentIntent(userID, orderID string) (clientSecret, paymentIntentID string, err error) {
	ord, err := o.GetOrder(userID, orderID)
	if err != nil {
		return "", "", err
	}
	if ord.PaymentStatus != "pending" || ord.Status != "created" {
		return "", "", ErrBadState
	}
	if strings.TrimSpace(o.StripeSecret) == "" {
		return "", "", ErrValidation
	}
	cur := strings.ToLower(strings.TrimSpace(o.StripeCurrency))
	if cur == "" {
		cur = "usd"
	}
	return payment.CreatePaymentIntent(o.StripeSecret, ord.Total, cur, ord.ID)
}

// QuoteShippingRates returns stub multi-carrier rates based on cart weight (FedEx/UPS/DHL style options).
func (o *OrderService) QuoteShippingRates(userID string, dest models.Address) ([]shipping.RateQuote, error) {
	cart, err := o.Cart.GetCart(userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, ErrValidation
	}
	var totalWeight float64
	for _, line := range cart.Items {
		p, err := o.Store.FindProductByID(line.ProductID)
		if err != nil || p == nil {
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
		w := v.WeightKg
		if w <= 0 {
			w = o.DefaultItemWeightKg
			if w <= 0 {
				w = 0.5
			}
		}
		totalWeight += w * float64(line.Quantity)
	}
	return shipping.QuoteCarriers(dest.Country, totalWeight), nil
}
