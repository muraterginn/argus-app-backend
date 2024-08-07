package main

import (
	"argus-app-backend/internal/config"
	"argus-app-backend/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	err = server.StartTCPServer(cfg)
	if err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}
