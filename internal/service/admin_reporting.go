package service

import (
	"sort"
	"strings"
	"time"
)

type OrderTimelineEvent struct {
	At    string `json:"at"`
	Label string `json:"label"`
	Kind  string `json:"kind"`
}

type DashboardStats struct {
	WindowDays   int     `json:"windowDays"`
	TotalRevenue float64 `json:"totalRevenue"`
	OrdersPlaced int     `json:"ordersPlaced"`
	PaidOrders   int     `json:"paidOrders"`
	NewCustomers int     `json:"newCustomers"`
	Previous     struct {
		TotalRevenue float64 `json:"totalRevenue"`
		OrdersPlaced int     `json:"ordersPlaced"`
		PaidOrders   int     `json:"paidOrders"`
		NewCustomers int     `json:"newCustomers"`
	} `json:"previous"`
}

func dashboardWindow(days int) (curStart, curEnd, prevStart, prevEnd time.Time) {
	if days <= 0 {
		days = 30
	}
	curEnd = time.Now().UTC()
	curStart = curEnd.Add(-time.Duration(days) * 24 * time.Hour)
	prevEnd = curStart
	prevStart = curStart.Add(-time.Duration(days) * 24 * time.Hour)
	return curStart, curEnd, prevStart, prevEnd
}

func timeInDashboardRange(t time.Time, start, end time.Time) bool {
	if t.IsZero() {
		return false
	}
	return !t.Before(start) && !t.After(end)
}

func parseOrderTime(ts string) time.Time {
	ts = strings.TrimSpace(ts)
	if ts == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return time.Time{}
	}
	return t
}

// AdminDashboardStats aggregates storefront KPIs for the last N days vs the prior N days.
func (o *OrderService) AdminDashboardStats(days int) (*DashboardStats, error) {
	curStart, curEnd, prevStart, prevEnd := dashboardWindow(days)
	if days <= 0 {
		days = 30
	}
	out := &DashboardStats{WindowDays: days}

	orders, err := o.Store.ListOrders()
	if err != nil {
		return nil, err
	}
	for _, ord := range orders {
		ct := parseOrderTime(ord.CreatedAt)
		paid := strings.EqualFold(ord.PaymentStatus, "paid") && !strings.EqualFold(ord.Status, "cancelled")
		if timeInDashboardRange(ct, curStart, curEnd) {
			out.OrdersPlaced++
			if paid {
				out.PaidOrders++
				out.TotalRevenue += ord.Total
			}
		}
		if timeInDashboardRange(ct, prevStart, prevEnd) {
			out.Previous.OrdersPlaced++
			if paid {
				out.Previous.PaidOrders++
				out.Previous.TotalRevenue += ord.Total
			}
		}
	}

	users, err := o.Store.ListUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		ut := parseOrderTime(u.CreatedAt)
		if timeInDashboardRange(ut, curStart, curEnd) {
			out.NewCustomers++
		}
		if timeInDashboardRange(ut, prevStart, prevEnd) {
			out.Previous.NewCustomers++
		}
	}
	return out, nil
}

// AdminOrderTimeline returns ordered lifecycle events for admin order detail.
func (o *OrderService) AdminOrderTimeline(orderID string) ([]OrderTimelineEvent, error) {
	ord, err := o.Store.FindOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord == nil {
		return nil, ErrNotFound
	}
	var ev []OrderTimelineEvent
	if ord.CreatedAt != "" {
		ev = append(ev, OrderTimelineEvent{At: ord.CreatedAt, Label: "Order placed", Kind: "created"})
	}
	paidAt := strings.TrimSpace(ord.PaidAt)
	pays, err := o.Store.ListPaymentsByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	for _, p := range pays {
		if strings.EqualFold(p.Status, "paid") {
			if paidAt == "" || parseOrderTime(p.CreatedAt).Before(parseOrderTime(paidAt)) {
				paidAt = p.CreatedAt
			}
		}
	}
	if paidAt != "" {
		ev = append(ev, OrderTimelineEvent{At: paidAt, Label: "Payment received", Kind: "paid"})
	}
	if strings.TrimSpace(ord.FulfilledAt) != "" {
		ev = append(ev, OrderTimelineEvent{At: ord.FulfilledAt, Label: "Order fulfilled", Kind: "fulfilled"})
	}
	if strings.TrimSpace(ord.ShippedAt) != "" {
		ev = append(ev, OrderTimelineEvent{At: ord.ShippedAt, Label: "Order shipped", Kind: "shipped"})
	}
	if strings.TrimSpace(ord.CancelledAt) != "" {
		ev = append(ev, OrderTimelineEvent{At: ord.CancelledAt, Label: "Order cancelled", Kind: "cancelled"})
	}
	sort.Slice(ev, func(i, j int) bool {
		ti := parseOrderTime(ev[i].At)
		tj := parseOrderTime(ev[j].At)
		if !ti.Equal(tj) {
			return ti.Before(tj)
		}
		return ev[i].Kind < ev[j].Kind
	})
	return ev, nil
}
