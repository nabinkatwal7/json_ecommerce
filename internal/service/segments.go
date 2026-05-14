package service

import (
	"strings"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"
)

// ComputeCustomerSegments derives marketing labels from paid order history.
func ComputeCustomerSegments(store *repository.Store, userID string, bigSpenderMin float64) ([]string, error) {
	if bigSpenderMin <= 0 {
		bigSpenderMin = 500
	}
	orders, err := store.OrdersByUser(userID)
	if err != nil {
		return nil, err
	}
	var paidCount int
	var total float64
	for _, o := range orders {
		if strings.EqualFold(o.Status, "cancelled") {
			continue
		}
		if o.PaymentStatus == "paid" {
			paidCount++
			total += o.Total
		}
	}
	var segs []string
	if paidCount == 0 {
		segs = append(segs, "first_time_buyer")
	} else {
		segs = append(segs, "returning_buyer")
	}
	if paidCount >= 2 {
		segs = append(segs, "repeat_buyer")
	}
	if total >= bigSpenderMin {
		segs = append(segs, "big_spender")
	}
	if paidCount >= 5 && total >= 2*bigSpenderMin {
		segs = append(segs, "vip")
	}
	return segs, nil
}

func (s *UserService) RefreshSegments(userID string, bigSpenderMin float64) (*models.User, error) {
	segs, err := ComputeCustomerSegments(s.Store, userID, bigSpenderMin)
	if err != nil {
		return nil, err
	}
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	u.Segments = segs
	if err := s.Store.UpsertUser(*u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}
