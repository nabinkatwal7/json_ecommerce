package api

import (
	"net/http"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) getCart(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	cart, err := r.Cart.GetCart(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (r *Router) postCartItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body service.AddCartItemInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	cart, err := r.Cart.AddItem(uid, body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}

type patchQtyBody struct {
	Quantity int `json:"quantity"`
}

func (r *Router) patchCartItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body patchQtyBody
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	cart, err := r.Cart.UpdateItemQty(uid, c.Param("itemId"), body.Quantity)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (r *Router) deleteCartItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	cart, err := r.Cart.RemoveItem(uid, c.Param("itemId"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}
