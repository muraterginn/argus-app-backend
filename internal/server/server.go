package server

import (
	"argus-app-backend/internal/config"
	"argus-app-backend/internal/tlsconfig"
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

func StartTLSServer(cfg config.Config) error {
	tlsConfig, err := tlsconfig.SetupTLSConfig(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", "127.0.0.1:"+cfg.ServerPort)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(ln, tlsConfig)
	defer tlsListener.Close()

	addr := tlsListener.Addr().String()
	log.Printf("TLS server started on %s", addr)

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("Received message: %s", message)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}
