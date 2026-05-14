package api

import (
	"net/http"
	"strconv"

	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) getSearch(c *gin.Context) {
	q := c.Query("q")
	cat := c.Query("categoryId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	hits, err := r.Search.Search(q, cat, limit)
	if err != nil {
		abortAPI(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"hits": hits})
}

func (r *Router) getSearchSuggest(c *gin.Context) {
	q := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	list, err := r.Discovery.SearchSuggestions(q, limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"suggestions": list})
}

func (r *Router) postAbandonedCartCron(c *gin.Context) {
	if r.Config.CronSecret == "" || c.GetHeader("X-Cron-Secret") != r.Config.CronSecret {
		abortAPI(c, http.StatusForbidden, "forbidden")
		return
	}
	n, err := r.Marketing.RunAbandonedCartEmails()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sent": n})
}

func (r *Router) postShippingQuote(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var addr models.Address
	if err := c.ShouldBindJSON(&addr); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	if addr.Country == "" {
		addr.Country = "US"
	}
	rates, err := r.Orders.QuoteShippingRates(uid, addr)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

func (r *Router) getMeInsights(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := r.Users.RefreshSegments(uid, r.Config.BigSpenderUSD)
	if err != nil {
		respondErr(c, err)
		return
	}
	orders, _ := r.Orders.ListMyOrders(uid)
	var spent float64
	paid := 0
	for _, o := range orders {
		if o.PaymentStatus == "paid" {
			paid++
			spent += o.Total
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"user":                 u,
		"paidOrders":           paid,
		"lifetimeSpend":        spent,
		"bigSpenderThresholdUsd": r.Config.BigSpenderUSD,
	})
}

func (r *Router) getWishlist(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := r.Lists.ListWishlist(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

type listItemBody struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
}

func (r *Router) postWishlistItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body listItemBody
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	list, err := r.Lists.AddWishlist(uid, body.ProductID, body.VariantID)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) deleteWishlistItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	pid := c.Query("productId")
	vid := c.Query("variantId")
	list, err := r.Lists.RemoveWishlist(uid, pid, vid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) postWishlistMoveSave(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body listItemBody
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	wish, later, err := r.Lists.MoveToSaveLater(uid, body.ProductID, body.VariantID)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"wishlist": wish, "saveLater": later})
}

func (r *Router) getSaveLater(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := r.Lists.ListSaveLater(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) postSaveLaterItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body listItemBody
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	list, err := r.Lists.AddSaveLater(uid, body.ProductID, body.VariantID)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) deleteSaveLaterItem(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := r.Lists.RemoveSaveLater(uid, c.Query("productId"), c.Query("variantId"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) postSaveLaterMoveWish(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body listItemBody
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	wish, later, err := r.Lists.MoveToWishlist(uid, body.ProductID, body.VariantID)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"wishlist": wish, "saveLater": later})
}

func (r *Router) postRMA(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body service.RMACreateInput
	if err := c.ShouldBindJSON(&body); err != nil {
		abortAPI(c, http.StatusBadRequest, "invalid json")
		return
	}
	rm, err := r.RMA.Create(uid, body)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, rm)
}

func (r *Router) getRMAs(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := r.RMA.ListMine(uid)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) getRMA(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	rm, err := r.RMA.GetMine(uid, c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rm)
}

type rmaNoteBody struct {
	Note string `json:"note"`
}

func (r *Router) adminListRMAs(c *gin.Context) {
	list, err := r.RMA.ListAdmin()
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (r *Router) adminGetRMA(c *gin.Context) {
	rm, err := r.RMA.GetAdmin(c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rm)
}

func (r *Router) adminApproveRMA(c *gin.Context) {
	var body rmaNoteBody
	_ = c.ShouldBindJSON(&body)
	rm, err := r.RMA.AdminApprove(c.Param("id"), body.Note)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rm)
}

func (r *Router) adminRejectRMA(c *gin.Context) {
	var body rmaNoteBody
	_ = c.ShouldBindJSON(&body)
	rm, err := r.RMA.AdminReject(c.Param("id"), body.Note)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rm)
}

func (r *Router) adminReceiveRMA(c *gin.Context) {
	var body rmaNoteBody
	_ = c.ShouldBindJSON(&body)
	rm, err := r.RMA.AdminMarkReceived(c.Param("id"), body.Note)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rm)
}

func (r *Router) adminRefundRMA(c *gin.Context) {
	var body rmaNoteBody
	_ = c.ShouldBindJSON(&body)
	rm, err := r.RMA.AdminRefund(c.Param("id"), body.Note)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rm)
}

func (r *Router) adminSearchReindex(c *gin.Context) {
	n, err := r.Search.ReindexAlgolia()
	if err != nil {
		abortAPI(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"indexed": n})
}
