package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds runtime settings for the API server and JSON store.
type Config struct {
	Addr       string
	DataDir    string
	JWTSecret  string
	Shipping   float64
	FreeShipAt float64

	TLSCertFile string
	TLSKeyFile  string

	StripeSecretKey   string
	StripeCurrency    string
	DevPaymentStub    bool
	RateLimitRPS      float64
	RateLimitBurst    int
	LoginRateLimitRPS float64
	LoginBurst        int

	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	AppPublicURL string

	LowStockThreshold int
	AdminAlertEmail   string

	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	CatalogCacheMem bool

	AlgoliaAppID  string
	AlgoliaAPIKey string
	AlgoliaIndex  string

	CronSecret                 string
	AbandonedCartMinHours      int
	AbandonedCartCooldownHours int

	DefaultItemWeightKg float64
	BigSpenderUSD       float64
}

func envBool(key string) bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(key)), "1") ||
		strings.EqualFold(strings.TrimSpace(os.Getenv(key)), "true") ||
		strings.EqualFold(strings.TrimSpace(os.Getenv(key)), "yes")
}

func envFloat(key string, def float64) float64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func envInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func Load() Config {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = filepath.Join(".", "data")
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-insecure-change-me"
	}
	cur := strings.TrimSpace(os.Getenv("STRIPE_CURRENCY"))
	if cur == "" {
		cur = "usd"
	}
	return Config{
		Addr:              addr,
		DataDir:           dataDir,
		JWTSecret:         secret,
		Shipping:          9.99,
		FreeShipAt:        50,
		TLSCertFile:       strings.TrimSpace(os.Getenv("TLS_CERT_FILE")),
		TLSKeyFile:        strings.TrimSpace(os.Getenv("TLS_KEY_FILE")),
		StripeSecretKey:   strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY")),
		StripeCurrency:    cur,
		DevPaymentStub:    envBool("DEV_PAYMENT_STUB"),
		RateLimitRPS:      envFloat("RATE_LIMIT_RPS", 25),
		RateLimitBurst:    envInt("RATE_LIMIT_BURST", 50),
		LoginRateLimitRPS: envFloat("LOGIN_RATE_LIMIT_RPS", 0.1),
		LoginBurst:        envInt("LOGIN_RATE_LIMIT_BURST", 5),
		SMTPHost:          strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:          strings.TrimSpace(os.Getenv("SMTP_PORT")),
		SMTPUser:          strings.TrimSpace(os.Getenv("SMTP_USER")),
		SMTPPassword:      strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		SMTPFrom:          strings.TrimSpace(os.Getenv("SMTP_FROM")),
		AppPublicURL:      strings.TrimSpace(os.Getenv("APP_PUBLIC_URL")),
		LowStockThreshold: envInt("LOW_STOCK_THRESHOLD", 5),
		AdminAlertEmail:   strings.TrimSpace(os.Getenv("ADMIN_ALERT_EMAIL")),

		RedisAddr:       strings.TrimSpace(os.Getenv("REDIS_ADDR")),
		RedisPassword:   strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
		RedisDB:         envInt("REDIS_DB", 0),
		CatalogCacheMem: envBool("CATALOG_CACHE_MEMORY"),

		AlgoliaAppID:  strings.TrimSpace(os.Getenv("ALGOLIA_APP_ID")),
		AlgoliaAPIKey: strings.TrimSpace(os.Getenv("ALGOLIA_API_KEY")),
		AlgoliaIndex:  strings.TrimSpace(os.Getenv("ALGOLIA_INDEX_NAME")),

		CronSecret:                 strings.TrimSpace(os.Getenv("CRON_SECRET")),
		AbandonedCartMinHours:      envInt("ABANDONED_CART_MIN_HOURS", 24),
		AbandonedCartCooldownHours: envInt("ABANDONED_CART_COOLDOWN_HOURS", 72),

		DefaultItemWeightKg: envFloat("DEFAULT_ITEM_WEIGHT_KG", 0.5),
		BigSpenderUSD:       envFloat("BIG_SPENDER_USD", 500),
	}
}
