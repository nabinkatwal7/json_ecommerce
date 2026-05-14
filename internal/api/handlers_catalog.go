package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getProducts(c *gin.Context) {
	categoryID := c.Query("categoryId")
	list, err := r.Catalog.ListActiveProducts(categoryID)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) getProduct(c *gin.Context) {
	p, err := r.Catalog.GetProduct(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (r *Router) getProductBySlug(c *gin.Context) {
	p, err := r.Catalog.GetProductBySlug(c.Param("slug"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (r *Router) getCategories(c *gin.Context) {
	list, err := r.Catalog.ListActiveCategories()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) getCategory(c *gin.Context) {
	cat, err := r.Catalog.GetCategory(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}
