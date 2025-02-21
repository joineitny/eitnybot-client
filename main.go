package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oschwald/geoip2-golang"
)

type Client struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Location string `json:"location"`
}

func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		log.Fatal("Error getting public IP:", err)
	}
	defer resp.Body.Close()

	ip := ""
	if _, err := fmt.Fscan(resp.Body, &ip); err != nil {
		log.Fatal("Error reading public IP:", err)
	}

	return ip
}

func getLocation(ip string) string {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal("Error opening GeoIP database:", err)
	}
	defer db.Close()

	record, err := db.City(net.ParseIP(ip))
	if err != nil {
		log.Fatal("Error getting location:", err)
	}

	return fmt.Sprintf("%s, %s", record.City.Names["en"], record.Country.Names["en"])
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Error getting hostname:", err)
	}
	return hostname
}

func main() {
	ip := getPublicIP()
	hostname := getHostname()
	location := getLocation(ip)

	client := Client{
		IP:       ip,
		Hostname: hostname,
		Location: location,
	}

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	for {
		clientJSON, err := json.Marshal(client)
		if err != nil {
			log.Println("Marshal error:", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, clientJSON); err != nil {
			log.Println("Write error:", err)
			break
		}

		time.Sleep(10 * time.Second)
	}
}
