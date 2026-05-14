package api

import (
	"net/http"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) adminListBanners(c *gin.Context) {
	list, err := r.Banners.AdminList()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) adminPostBanner(c *gin.Context) {
	var body service.BannerInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	bn, err := r.Banners.AdminCreate(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, bn)
}

func (r *Router) adminPutBanner(c *gin.Context) {
	var body service.BannerInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	bn, err := r.Banners.AdminUpdate(c.Param("id"), body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, bn)
}

func (r *Router) adminDeleteBanner(c *gin.Context) {
	if err := r.Banners.AdminDelete(c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
