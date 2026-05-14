package api

import (
	"net/http"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) getTags(c *gin.Context) {
	list, err := r.Tags.List()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) adminListTags(c *gin.Context) {
	list, err := r.Tags.List()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) adminPostTag(c *gin.Context) {
	var body service.TagInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	t, err := r.Tags.AdminCreate(body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (r *Router) adminPutTag(c *gin.Context) {
	var body service.TagInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	t, err := r.Tags.AdminUpdate(c.Param("id"), body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, t)
}

func (r *Router) adminDeleteTag(c *gin.Context) {
	if err := r.Tags.AdminDelete(c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
