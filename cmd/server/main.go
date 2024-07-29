package main

import (
	"argus-app-backend/internal/config"
	"argus-app-backend/internal/server"
	"log"
)

func main() {
	// Konfigürasyonu yükle
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// TLS sunucusunu başlat
	err = server.StartTLSServer(cfg)
	if err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}
