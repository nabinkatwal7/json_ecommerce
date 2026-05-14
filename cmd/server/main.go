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

	log.Printf("listening on %s (data: %s)", cfg.Addr, cfg.DataDir)
	if err := engine.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
