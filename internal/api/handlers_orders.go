package api

import (
	"net/http"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) postCheckout(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body service.CheckoutInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	order, err := r.Orders.Checkout(uid, body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (r *Router) getOrders(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := r.Orders.ListMyOrders(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) getOrder(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	o, err := r.Orders.GetOrder(uid, c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (r *Router) getOrderInvoicePDF(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	pdf, err := r.Orders.InvoicePDF(uid, c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Data(http.StatusOK, "application/pdf", pdf)
}

func (r *Router) postCancelOrder(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	o, err := r.Orders.CancelByCustomer(uid, c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (r *Router) postStripePaymentIntent(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	cs, id, err := r.Orders.CreateStripePaymentIntent(uid, c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"clientSecret": cs, "paymentIntentId": id})
}

func (r *Router) postPayOrder(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var in service.PayInput
	if err := c.ShouldBindJSON(&in); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	o, pay, err := r.Orders.Pay(uid, c.Param("id"), in)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": o, "payment": pay})
}
