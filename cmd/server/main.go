package main

import (
	"log"
	"os"

	"go-ecommerce-json/internal/api"
	"go-ecommerce-json/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	api.NewRouter(cfg).Mount(engine)

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		log.Printf("listening with TLS on %s (data: %s)", cfg.Addr, cfg.DataDir)
		if err := engine.RunTLS(cfg.Addr, cfg.TLSCertFile, cfg.TLSKeyFile); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Printf("listening on %s (data: %s) — set TLS_CERT_FILE and TLS_KEY_FILE for HTTPS", cfg.Addr, cfg.DataDir)
	if err := engine.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
