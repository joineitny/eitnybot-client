package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	if len(os.Args) < 6 {
		fmt.Println("Использование: client <ID> <IP> <Port> <Username> <Password>")
		return
	}

	comp := Computer{
		ID:       os.Args[1],
		IP:       os.Args[2],
		Port:     os.Args[3],
		Username: os.Args[4],
		Password: os.Args[5],
	}

	registerComputer(comp)
}

func registerComputer(comp Computer) {
	jsonData, err := json.Marshal(comp)
	if err != nil {
		log.Fatalf("Ошибка при сериализации данных: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080/register", bytes.NewBuffer(jsonData))
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

	log.Println("Компьютер успешно зарегистрирован на сервере")
}
