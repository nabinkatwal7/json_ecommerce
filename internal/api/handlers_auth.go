package api

import (
	"net/http"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) postRegister(c *gin.Context) {
	var body service.RegisterInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	res, err := r.Users.Register(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (r *Router) postLogin(c *gin.Context) {
	var body service.LoginInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	res, err := r.Users.Login(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (r *Router) getMe(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := r.Users.GetProfile(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}
