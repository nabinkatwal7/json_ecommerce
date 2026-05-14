package config

import (
	"os"
	"path/filepath"
)

// Config holds runtime settings for the API server and JSON store.
type Config struct {
	Addr       string
	DataDir    string
	JWTSecret  string
	Shipping   float64
	FreeShipAt float64
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
	return Config{
		Addr:       addr,
		DataDir:    dataDir,
		JWTSecret:  secret,
		Shipping:   9.99,
		FreeShipAt: 50,
	}
}
