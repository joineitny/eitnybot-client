package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type ClientInfo struct {
	IP       string
	Hostname string
}

func main() {
	hostname, _ := os.Hostname()
	info := ClientInfo{
		IP:       getOutboundIP(),
		Hostname: hostname,
	}

	conn, err := net.Dial("tcp", "server-ip:8080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(info); err != nil {
		log.Fatalf("Failed to send client info: %v", err)
	}

	// Keep the connection alive
	for {
		time.Sleep(10 * time.Second)
	}
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
