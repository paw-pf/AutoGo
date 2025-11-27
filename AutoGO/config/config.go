package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

var (
	APIBaseURL string
	UIBaseURL  string
	Username   string
	Password   string
	Token      string
	Headless   bool
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	envPath := filepath.Join(basepath, "..", ".env")

	log.Printf("🔍 Ищу .env по пути: %s", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  Не удалось загрузить .env: %v", err)
	} else {
		log.Println("✅ .env успешно загружен")
	}

	APIBaseURL = os.Getenv("API_BASE_URL")
	UIBaseURL = os.Getenv("UI_BASE_URL")
	Username = os.Getenv("USERNAME")
	Password = os.Getenv("PASSWORD")
	Token = os.Getenv("TOKEN")
	Headless = getEnvBool("HEADLESS", true)

}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return fallback
}
