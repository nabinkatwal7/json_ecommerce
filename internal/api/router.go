package api

import (
	"errors"
	"net/http"
	"time"

	"go-ecommerce-json/internal/auth"
	"go-ecommerce-json/internal/cache"
	"go-ecommerce-json/internal/config"
	"go-ecommerce-json/internal/mail"
	"go-ecommerce-json/internal/repository"
	"go-ecommerce-json/internal/search"
	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "userID"
	ctxRole   = "role"
)

type Router struct {
	Config    config.Config
	Store     *repository.Store
	Users     *service.UserService
	Catalog   *service.CatalogService
	Cart      *service.CartService
	Orders    *service.OrderService
	Promo     *service.PromoService
	Tags      *service.TagService
	Password  *service.PasswordResetService
	Marketing *service.MarketingService
	RMA       *service.RMAService
	Lists     *service.UserListsService
	Search    *search.Client

	defaultLim *ipLimiter
	loginLim   *ipLimiter
}

func NewRouter(cfg config.Config) *Router {
	st := repository.NewStore(cfg.DataDir)
	m := &mail.Sender{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	}
	us := &service.UserService{
		Store:     st,
		JWTSecret: []byte(cfg.JWTSecret),
		JWTTTL:    7 * 24 * time.Hour,
	}

	var catCache cache.CatalogCache = cache.NoOp()
	if cfg.RedisAddr != "" {
		if rdb, err := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB); err == nil {
			catCache = cache.NewRedisCatalogCache(rdb)
		}
	} else if cfg.CatalogCacheMem {
		catCache = cache.NewMemoryCatalogCache()
	}

	cs := &service.CatalogService{
		Store:    st,
		Cache:    catCache,
		CacheTTL: 30 * time.Second,
	}
	cart := &service.CartService{Store: st}
	orders := &service.OrderService{
		Store:               st,
		Cart:                cart,
		Shipping:            cfg.Shipping,
		FreeShipAt:          cfg.FreeShipAt,
		DefaultItemWeightKg: cfg.DefaultItemWeightKg,
		Mail:                m,
		StripeSecret:        cfg.StripeSecretKey,
		StripeCurrency:      cfg.StripeCurrency,
		DevPaymentStub:      cfg.DevPaymentStub,
		AppPublicURL:        cfg.AppPublicURL,
		LowStockThreshold:   cfg.LowStockThreshold,
		AdminAlertEmail:     cfg.AdminAlertEmail,
	}
	promo := &service.PromoService{Store: st}
	tags := &service.TagService{Store: st}
	pw := &service.PasswordResetService{
		Store:    st,
		Mail:     m,
		AppURL:   cfg.AppPublicURL,
		TokenTTL: time.Hour,
	}
	mrk := &service.MarketingService{
		Store:    st,
		Mail:     m,
		AppURL:   cfg.AppPublicURL,
		MinIdle:  time.Duration(cfg.AbandonedCartMinHours) * time.Hour,
		Cooldown: time.Duration(cfg.AbandonedCartCooldownHours) * time.Hour,
	}
	rma := &service.RMAService{Store: st}
	lists := &service.UserListsService{Store: st}
	srch := &search.Client{
		Store:         st,
		AlgoliaAppID:  cfg.AlgoliaAppID,
		AlgoliaAPIKey: cfg.AlgoliaAPIKey,
		AlgoliaIndex:  cfg.AlgoliaIndex,
	}
	return &Router{
		Config:     cfg,
		Store:      st,
		Users:      us,
		Catalog:    cs,
		Cart:       cart,
		Orders:     orders,
		Promo:      promo,
		Tags:       tags,
		Password:   pw,
		Marketing:  mrk,
		RMA:        rma,
		Lists:      lists,
		Search:     srch,
		defaultLim: newIPLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst),
		loginLim:   newIPLimiter(cfg.LoginRateLimitRPS, cfg.LoginBurst),
	}
}

func (r *Router) Mount(engine *gin.Engine) {
	engine.Use(r.defaultLim.middleware())

	engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	engine.POST("/register", r.loginLim.middleware(), r.postRegister)
	engine.POST("/login", r.loginLim.middleware(), r.postLogin)
	engine.POST("/forgot-password", r.loginLim.middleware(), r.postForgotPassword)
	engine.POST("/reset-password", r.loginLim.middleware(), r.postResetPassword)

	engine.GET("/products", r.getProducts)
	engine.GET("/products/:id", r.getProduct)
	engine.GET("/products/slug/:slug", r.getProductBySlug)
	engine.GET("/categories", r.getCategories)
	engine.GET("/categories/:id", r.getCategory)
	engine.GET("/tags", r.getTags)
	engine.GET("/search", r.getSearch)

	engine.POST("/internal/cron/abandoned-carts", r.loginLim.middleware(), r.postAbandonedCartCron)

	authz := engine.Group("/")
	authz.Use(r.authMiddleware())
	authz.GET("/me", r.getMe)
	authz.GET("/me/insights", r.getMeInsights)

	authz.GET("/me/addresses", r.getAddresses)
	authz.POST("/me/addresses", r.postAddress)
	authz.PUT("/me/addresses/:id", r.putAddress)
	authz.DELETE("/me/addresses/:id", r.deleteAddress)

	authz.POST("/shipping/quote", r.postShippingQuote)

	authz.GET("/wishlist", r.getWishlist)
	authz.POST("/wishlist/items", r.postWishlistItem)
	authz.DELETE("/wishlist/items", r.deleteWishlistItem)
	authz.POST("/wishlist/move-to-save-later", r.postWishlistMoveSave)

	authz.GET("/save-later", r.getSaveLater)
	authz.POST("/save-later/items", r.postSaveLaterItem)
	authz.DELETE("/save-later/items", r.deleteSaveLaterItem)
	authz.POST("/save-later/move-to-wishlist", r.postSaveLaterMoveWish)

	authz.GET("/cart", r.getCart)
	authz.POST("/cart/items", r.postCartItem)
	authz.PATCH("/cart/items/:itemId", r.patchCartItem)
	authz.DELETE("/cart/items/:itemId", r.deleteCartItem)

	authz.POST("/rmas", r.postRMA)
	authz.GET("/rmas", r.getRMAs)
	authz.GET("/rmas/:id", r.getRMA)

	authz.POST("/orders/checkout", r.postCheckout)
	authz.GET("/orders", r.getOrders)
	authz.GET("/orders/:id/invoice.pdf", r.getOrderInvoicePDF)
	authz.POST("/orders/:id/cancel", r.postCancelOrder)
	authz.POST("/orders/:id/stripe-payment-intent", r.postStripePaymentIntent)
	authz.POST("/orders/:id/pay", r.postPayOrder)
	authz.GET("/orders/:id", r.getOrder)

	admin := engine.Group("/admin")
	admin.Use(r.authMiddleware(), r.adminMiddleware())
	admin.GET("/products", r.adminListProducts)
	admin.POST("/products", r.adminPostProduct)
	admin.PUT("/products/:id", r.adminPutProduct)
	admin.DELETE("/products/:id", r.adminDeleteProduct)

	admin.POST("/categories", r.adminPostCategory)
	admin.PUT("/categories/:id", r.adminPutCategory)
	admin.DELETE("/categories/:id", r.adminDeleteCategory)

	admin.POST("/discounts", r.adminPostDiscount)

	admin.GET("/tags", r.adminListTags)
	admin.POST("/tags", r.adminPostTag)
	admin.PUT("/tags/:id", r.adminPutTag)
	admin.DELETE("/tags/:id", r.adminDeleteTag)

	admin.GET("/orders", r.adminListOrders)
	admin.POST("/orders/:id/cancel", r.adminCancelOrder)
	admin.POST("/orders/:id/fulfill", r.adminFulfillOrder)
	admin.POST("/orders/:id/ship", r.adminShipOrder)

	admin.GET("/inventory/low-stock", r.adminLowStock)

	admin.GET("/rmas", r.adminListRMAs)
	admin.GET("/rmas/:id", r.adminGetRMA)
	admin.POST("/rmas/:id/approve", r.adminApproveRMA)
	admin.POST("/rmas/:id/reject", r.adminRejectRMA)
	admin.POST("/rmas/:id/receive", r.adminReceiveRMA)
	admin.POST("/rmas/:id/refund", r.adminRefundRMA)

	admin.POST("/search/reindex", r.adminSearchReindex)
}

func (r *Router) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if len(h) <= len(prefix) || h[:len(prefix)] != prefix {
			abortAPI(c, http.StatusUnauthorized, "missing bearer token")
			return
		}
		token := h[len(prefix):]
		claims, err := auth.ParseJWT([]byte(r.Config.JWTSecret), token)
		if err != nil {
			abortAPI(c, http.StatusUnauthorized, "invalid token")
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

func (r *Router) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ctxRole)
		rs, _ := role.(string)
		if rs != "admin" {
			abortAPI(c, http.StatusForbidden, "admin only")
			return
		}
		c.Next()
	}
}

func userID(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}

func abortAPI(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{"error": msg})
}

func respondErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		abortAPI(c, http.StatusNotFound, "not found")
	case errors.Is(err, service.ErrConflict):
		abortAPI(c, http.StatusConflict, "conflict")
	case errors.Is(err, service.ErrUnauthorized):
		abortAPI(c, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, service.ErrForbidden):
		abortAPI(c, http.StatusForbidden, "forbidden")
	case errors.Is(err, service.ErrValidation):
		abortAPI(c, http.StatusBadRequest, "validation failed")
	case errors.Is(err, service.ErrInsufficientStock):
		abortAPI(c, http.StatusBadRequest, "insufficient stock")
	case errors.Is(err, service.ErrInactive):
		abortAPI(c, http.StatusBadRequest, "inactive")
	case errors.Is(err, service.ErrBadState):
		abortAPI(c, http.StatusConflict, "invalid state")
	case errors.Is(err, service.ErrPayment):
		abortAPI(c, http.StatusPaymentRequired, "payment verification failed")
	default:
		abortAPI(c, http.StatusInternalServerError, "internal error")
	}
}
