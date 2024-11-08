package server

import (
	"argus-app-backend/internal/config"
	"argus-app-backend/internal/tlsconfig"
	"argus-app-backend/internal/utils"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]bool)
var mu sync.Mutex

func StartTCPServer(cfg config.Config) error {
	tlsConfig, err := tlsconfig.SetupTLSConfig(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", GetLocalIPv4Address()+":"+cfg.TCPServerPort)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(ln, tlsConfig)
	defer tlsListener.Close()

	addr := tlsListener.Addr().String()
	log.Printf("TCP server started on %s and waiting for connections...", addr)

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		ip := strings.Split(remoteAddr, ":")[0]
		if !isAllowedAddress(ip, cfg) {
			log.Printf("Connection denied from %s", ip)
			conn.Close()
			continue
		}

		log.Printf("Connection accepted from %s", ip)

		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		go handleConnection(conn)
	}
}

func GetPublicIPv4Address() string {
	response, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "Error fetching public IP"
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "Error reading response body"
	}

	return string(body)
}

func GetLocalIPv4Address() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error getting local IP addresses: %v", err)
		return "Error fetching local IP"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}

	return "Local IP not found"
}

func handleConnection(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	buffer := make([]byte, 16384)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			break
		}
		message := string(buffer[:n])

		if strings.HasPrefix(message, "[ACCOUNT]") {
			jsonData, err := utils.ParseMessageToJSON(message)
			if err != nil {
				log.Printf("Error parsing message to JSON: %v", err)
				continue
			}
			message = jsonData
		}

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

func isAllowedAddress(ip string, cfg config.Config) bool {
	if len(cfg.AllowedAddresses) == 0 {
		log.Println("No allowed addresses specified. Connection denied by default.")
		return false
	}
	for _, allowed := range cfg.AllowedAddresses {
		fmt.Print(allowed)
		if strings.TrimSpace(allowed) == ip {
			return true
		}
	}
	log.Printf("IP %s is not in the allowed addresses list. Connection denied.", ip)
	return false
}
