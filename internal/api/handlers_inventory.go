package api

import (
	"net/http"
	"strconv"

	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) adminLowStock(c *gin.Context) {
	th := r.Config.LowStockThreshold
	if q := c.Query("threshold"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			th = v
		}
	}
	lines, err := service.LowStockReport(r.Store, th)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"threshold": th, "lines": lines})
}
