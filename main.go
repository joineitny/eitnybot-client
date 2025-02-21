package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type Computer struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Автоматическое определение IP-адреса
	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Не удалось определить IP-адрес: %v", err)
	}

	comp := Computer{
		IP: ip,
	}

	registerComputer(comp)
}

// Функция для получения локального IP-адреса
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("не удалось определить IP-адрес")
}

func registerComputer(comp Computer) {
	jsonData, err := json.Marshal(comp)
	if err != nil {
		log.Fatalf("Ошибка при сериализации данных: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Указываем IP-адрес сервера 37.46.230.242
	serverURL := "http://37.46.230.242:8080/register"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Сервер вернул ошибку: %s, тело ответа: %s", resp.Status, string(body))
	}

	// Чтение ответа от сервера
	var response struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Ошибка при чтении ответа от сервера: %v", err)
	}

	log.Printf("Ответ от сервера: %s\n", response.Message)
}
