package service

import (
	"time"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type CartService struct {
	Store *repository.Store
}

func (c *CartService) getOrCreateCart(userID string) (*models.Cart, error) {
	cart, err := c.Store.FindCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart != nil {
		return cart, nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	nc := models.Cart{
		ID:        uuid.NewString(),
		UserID:    userID,
		Items:     []models.CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := c.Store.UpsertCart(nc); err != nil {
		return nil, err
	}
	return &nc, nil
}

func (c *CartService) GetCart(userID string) (*models.Cart, error) {
	return c.getOrCreateCart(userID)
}

type AddCartItemInput struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
	Quantity  int    `json:"quantity"`
}

func (c *CartService) AddItem(userID string, in AddCartItemInput) (*models.Cart, error) {
	if in.Quantity <= 0 {
		return nil, ErrValidation
	}
	p, err := c.Store.FindProductByID(in.ProductID)
	if err != nil {
		return nil, err
	}
	if p == nil || !p.IsActive {
		return nil, ErrNotFound
	}
	var v *models.ProductVariant
	for i := range p.Variants {
		if p.Variants[i].ID == in.VariantID {
			v = &p.Variants[i]
			break
		}
	}
	if v == nil {
		return nil, ErrValidation
	}
	if v.Stock < in.Quantity {
		return nil, ErrInsufficientStock
	}
	cart, err := c.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	found := false
	for i := range cart.Items {
		if cart.Items[i].ProductID == in.ProductID && cart.Items[i].VariantID == in.VariantID {
			newQty := cart.Items[i].Quantity + in.Quantity
			if newQty > v.Stock {
				return nil, ErrInsufficientStock
			}
			cart.Items[i].Quantity = newQty
			cart.Items[i].Price = v.Price
			cart.Items[i].Name = p.Name
			cart.Items[i].SKU = v.SKU
			cart.Items[i].Image = p.Image
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, models.CartItem{
			ID:        uuid.NewString(),
			ProductID: p.ID,
			VariantID: v.ID,
			Name:      p.Name,
			SKU:       v.SKU,
			Price:     v.Price,
			Quantity:  in.Quantity,
			Image:     p.Image,
		})
	}
	cart.UpdatedAt = now
	if err := c.Store.UpsertCart(*cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (c *CartService) UpdateItemQty(userID, itemID string, quantity int) (*models.Cart, error) {
	if quantity <= 0 {
		return nil, ErrValidation
	}
	cart, err := c.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	idx := -1
	for i := range cart.Items {
		if cart.Items[i].ID == itemID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, ErrNotFound
	}
	line := cart.Items[idx]
	p, err := c.Store.FindProductByID(line.ProductID)
	if err != nil {
		return nil, err
	}
	if p == nil || !p.IsActive {
		return nil, ErrNotFound
	}
	var v *models.ProductVariant
	for i := range p.Variants {
		if p.Variants[i].ID == line.VariantID {
			v = &p.Variants[i]
			break
		}
	}
	if v == nil || v.Stock < quantity {
		return nil, ErrInsufficientStock
	}
	cart.Items[idx].Quantity = quantity
	cart.Items[idx].Price = v.Price
	cart.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := c.Store.UpsertCart(*cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (c *CartService) RemoveItem(userID, itemID string) (*models.Cart, error) {
	cart, err := c.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	out := cart.Items[:0]
	removed := false
	for _, it := range cart.Items {
		if it.ID == itemID {
			removed = true
			continue
		}
		out = append(out, it)
	}
	if !removed {
		return nil, ErrNotFound
	}
	cart.Items = out
	cart.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := c.Store.UpsertCart(*cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (c *CartService) ClearCart(cartID string) error {
	carts, err := c.Store.ListCarts()
	if err != nil {
		return err
	}
	for i := range carts {
		if carts[i].ID == cartID {
			carts[i].Items = []models.CartItem{}
			carts[i].UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			return c.Store.SaveCarts(carts)
		}
	}
	return ErrNotFound
}

// ClearCartByUser clears the cart belonging to the user (used after successful payment).
func (c *CartService) ClearCartByUser(userID string) error {
	cart, err := c.Store.FindCartByUserID(userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}
	return c.ClearCart(cart.ID)
}
