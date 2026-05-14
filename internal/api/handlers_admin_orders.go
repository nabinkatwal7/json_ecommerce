package api

import (
	"net/http"

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
