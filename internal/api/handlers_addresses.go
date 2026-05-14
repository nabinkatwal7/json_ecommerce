package api

import (
	"net/http"

	"go-ecommerce-json/internal/models"

	"github.com/gin-gonic/gin"
)

func (r *Router) getAddresses(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := r.Users.ListAddresses(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) postAddress(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body models.Address
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := r.Users.AddAddress(uid, body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (r *Router) putAddress(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body models.Address
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := r.Users.UpdateAddress(uid, c.Param("id"), body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func (r *Router) deleteAddress(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := r.Users.DeleteAddress(uid, c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}
