package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Computer struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`       // Внешний IP-адрес
	Port     string `json:"port"`     // Порт для подключения
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль
}

func main() {
	// Автоматическое определение внешнего IP-адреса
	ip, err := getExternalIP()
	if err != nil {
		log.Fatalf("Не удалось определить внешний IP-адрес: %v", err)
	}

	comp := Computer{
		IP:   ip,
		Port: "8081", // Указываем порт, который будет использоваться для подключения
	}

	registerComputer(comp)
}

// Функция для получения внешнего IP-адреса
func getExternalIP() (string, error) {
	// Используем внешний сервис для получения IP-адреса
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("ошибка при получении внешнего IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("сервер вернул ошибку: %s", resp.Status)
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	return string(ip), nil
}

// Регистрация компьютера на сервере
func registerComputer(comp Computer) {
	jsonData, err := json.Marshal(comp)
	if err != nil {
		log.Fatalf("Ошибка при сериализации данных: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Указываем IP-адрес сервера
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
