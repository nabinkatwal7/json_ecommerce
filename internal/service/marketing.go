package service

import (
	"fmt"
	"strings"
	"time"

	"go-ecommerce-json/internal/mail"
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"
)

// MarketingService sends lifecycle emails (e.g. abandoned cart recovery).
type MarketingService struct {
	Store    *repository.Store
	Mail     *mail.Sender
	AppURL   string
	MinIdle  time.Duration
	Cooldown time.Duration
}

func (m *MarketingService) cartEligible(c models.Cart, now time.Time) bool {
	if len(c.Items) == 0 {
		return false
	}
	t, err := time.Parse(time.RFC3339, c.UpdatedAt)
	if err != nil {
		return false
	}
	if now.Sub(t) < m.MinIdle {
		return false
	}
	if strings.TrimSpace(c.LastAbandonedEmailAt) != "" {
		last, err := time.Parse(time.RFC3339, c.LastAbandonedEmailAt)
		if err == nil && now.Sub(last) < m.Cooldown {
			return false
		}
	}
	return true
}

// RunAbandonedCartEmails notifies users with non-empty carts idle for MinIdle and not reminded within Cooldown.
func (m *MarketingService) RunAbandonedCartEmails() (sent int, err error) {
	if m.Mail == nil || !m.Mail.Configured() {
		return 0, nil
	}
	carts, err := m.Store.ListCarts()
	if err != nil {
		return 0, err
	}
	now := time.Now().UTC()
	for i := range carts {
		c := carts[i]
		if !m.cartEligible(c, now) {
			continue
		}
		u, err := m.Store.FindUserByID(c.UserID)
		if err != nil || u == nil || strings.TrimSpace(u.Email) == "" {
			continue
		}
		link := strings.TrimRight(m.AppURL, "/") + "/cart"
		if m.AppURL == "" {
			link = "your cart"
		}
		var lines []string
		for _, it := range c.Items {
			lines = append(lines, fmt.Sprintf("- %s x%d (%s)", it.Name, it.Quantity, it.SKU))
		}
		body := "You still have items in your cart:\n\n" + strings.Join(lines, "\n") + "\n\nContinue checkout: " + link + "\n"
		if err := m.Mail.SendPlain(u.Email, "Complete your order", body); err != nil {
			continue
		}
		c.LastAbandonedEmailAt = now.Format(time.RFC3339)
		carts[i] = c
		if err := m.Store.UpsertCart(c); err != nil {
			return sent, err
		}
		sent++
	}
	return sent, nil
}
