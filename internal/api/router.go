package api

import (
	"errors"
	"net/http"
	"time"

	"go-ecommerce-json/internal/auth"
	"go-ecommerce-json/internal/config"
	"go-ecommerce-json/internal/repository"
	"go-ecommerce-json/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "userID"
	ctxRole   = "role"
)

type Router struct {
	Config  config.Config
	Store   *repository.Store
	Users   *service.UserService
	Catalog *service.CatalogService
	Cart    *service.CartService
	Orders  *service.OrderService
	Promo   *service.PromoService
}

func NewRouter(cfg config.Config) *Router {
	st := repository.NewStore(cfg.DataDir)
	us := &service.UserService{
		Store:     st,
		JWTSecret: []byte(cfg.JWTSecret),
		JWTTTL:    7 * 24 * time.Hour,
	}
	cs := &service.CatalogService{Store: st}
	cart := &service.CartService{Store: st}
	orders := &service.OrderService{
		Store:      st,
		Cart:       cart,
		Shipping:   cfg.Shipping,
		FreeShipAt: cfg.FreeShipAt,
	}
	promo := &service.PromoService{Store: st}
	return &Router{
		Config:  cfg,
		Store:   st,
		Users:   us,
		Catalog: cs,
		Cart:    cart,
		Orders:  orders,
		Promo:   promo,
	}
}

func (r *Router) Mount(engine *gin.Engine) {
	engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	engine.POST("/register", r.postRegister)
	engine.POST("/login", r.postLogin)

	pub := engine.Group("/")
	pub.GET("/products", r.getProducts)
	pub.GET("/products/:id", r.getProduct)
	pub.GET("/products/slug/:slug", r.getProductBySlug)
	pub.GET("/categories", r.getCategories)

	authz := engine.Group("/")
	authz.Use(r.authMiddleware())
	authz.GET("/me", r.getMe)

	authz.GET("/cart", r.getCart)
	authz.POST("/cart/items", r.postCartItem)
	authz.PATCH("/cart/items/:itemId", r.patchCartItem)
	authz.DELETE("/cart/items/:itemId", r.deleteCartItem)

	authz.POST("/orders/checkout", r.postCheckout)
	authz.GET("/orders", r.getOrders)
	authz.GET("/orders/:id", r.getOrder)
	authz.POST("/orders/:id/pay", r.postPayOrder)

	admin := engine.Group("/admin")
	admin.Use(r.authMiddleware(), r.adminMiddleware())
	admin.POST("/products", r.adminPostProduct)
	admin.PUT("/products/:id", r.adminPutProduct)
	admin.DELETE("/products/:id", r.adminDeleteProduct)
	admin.POST("/categories", r.adminPostCategory)
	admin.POST("/discounts", r.adminPostDiscount)
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
	default:
		abortAPI(c, http.StatusInternalServerError, "internal error")
	}
}
