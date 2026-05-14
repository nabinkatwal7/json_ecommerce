package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) adminListOrders(c *gin.Context) {
	list, err := r.Orders.ListOrdersAdmin()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) adminCancelOrder(c *gin.Context) {
	o, err := r.Orders.AdminCancel(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (r *Router) adminFulfillOrder(c *gin.Context) {
	o, err := r.Orders.AdminFulfill(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (r *Router) adminShipOrder(c *gin.Context) {
	o, err := r.Orders.AdminShip(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (r *Router) adminDashboardStats(c *gin.Context) {
	days := 30
	if q := c.Query("days"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 366 {
			days = n
		}
	}
	out, err := r.Orders.AdminDashboardStats(days)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (r *Router) adminOrderTimeline(c *gin.Context) {
	ev, err := r.Orders.AdminOrderTimeline(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": ev})
}
