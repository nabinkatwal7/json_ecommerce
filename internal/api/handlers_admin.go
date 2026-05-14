package api

import (
	"net/http"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) adminListProducts(c *gin.Context) {
	list, err := r.Catalog.AdminListProducts()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) adminPostProduct(c *gin.Context) {
	var body service.ProductInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	p, err := r.Catalog.AdminCreateProduct(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (r *Router) adminPutProduct(c *gin.Context) {
	var body service.ProductInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	p, err := r.Catalog.AdminUpdateProduct(c.Param("id"), body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (r *Router) adminDeleteProduct(c *gin.Context) {
	if err := r.Catalog.AdminDeleteProduct(c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (r *Router) adminPostCategory(c *gin.Context) {
	var body service.CategoryInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	cat, err := r.Catalog.AdminCreateCategory(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (r *Router) adminPutCategory(c *gin.Context) {
	var body service.CategoryInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	cat, err := r.Catalog.AdminUpdateCategory(c.Param("id"), body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (r *Router) adminDeleteCategory(c *gin.Context) {
	if err := r.Catalog.AdminDeleteCategory(c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (r *Router) adminPostDiscount(c *gin.Context) {
	var body service.DiscountInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	d, err := r.Promo.AdminCreateDiscount(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, d)
}
