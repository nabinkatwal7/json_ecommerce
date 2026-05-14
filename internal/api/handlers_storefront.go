package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) getFeaturedCollections(c *gin.Context) {
	out, err := r.Discovery.FeaturedHome()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (r *Router) getBanners(c *gin.Context) {
	slot := c.Query("slot")
	list, err := r.Banners.ListPublic(slot)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) getStorefrontFeed(c *gin.Context) {
	feat, err := r.Discovery.FeaturedHome()
	if err != nil {
		respondErr(c, err)
		return
	}
	sale, err := r.Discovery.SaleProducts(10)
	if err != nil {
		respondErr(c, err)
		return
	}
	ann, err := r.Banners.ListPublic("announcement")
	if err != nil {
		respondErr(c, err)
		return
	}
	car, err := r.Banners.ListPublic("home_carousel")
	if err != nil {
		respondErr(c, err)
		return
	}
	pc, cc, err := r.Discovery.CatalogCounts()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"featured":            feat,
		"salePicks":           sale,
		"announcementBanners": ann,
		"carouselBanners":     car,
		"saleTagId":           r.Discovery.TagIDBySlug("sale"),
		"stats": gin.H{
			"productCount":      pc,
			"categoryCount":     cc,
			"freeShippingAtUsd": r.Config.FreeShipAt,
		},
	})
}

func (r *Router) getProductRelated(c *gin.Context) {
	limit := 8
	if q := c.Query("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			limit = n
		}
	}
	list, err := r.Discovery.RelatedProducts(c.Param("id"), limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
