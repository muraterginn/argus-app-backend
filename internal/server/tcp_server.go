package server

import (
	"argus-app-backend/internal/config"
	"argus-app-backend/internal/tlsconfig"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
)

var clients = make(map[net.Conn]bool)
var mu sync.Mutex

func StartTCPServer(cfg config.Config) error {
	tlsConfig, err := tlsconfig.SetupTLSConfig(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", "127.0.0.1:"+cfg.TCPServerPort)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(ln, tlsConfig)
	defer tlsListener.Close()

	addr := tlsListener.Addr().String()
	log.Printf("TCP server started on %s", addr)

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			break
		}
		message := string(buffer[:n])
		log.Printf("%s", message)
		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		if client != sender {
			_, err := client.Write([]byte(message))
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
			}
		}
	}
}
