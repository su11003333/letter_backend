// backend/configs/config.go
package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 應用程式配置
type Config struct {
	Port             string
	JWTSecret        []byte
	JWTExpiryTime    time.Duration
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// LoadConfig 從環境變數載入配置
func LoadConfig() *Config {
	// 載入 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 設置默認值
	config := &Config{
		Port:             getEnv("PORT", "8080"),
		JWTSecret:        []byte(getEnv("JWT_SECRET", "your_secure_secret_key")),
		JWTExpiryTime:    time.Duration(getEnvAsInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,
		AllowedOrigins:   []string{getEnv("ALLOWED_ORIGINS", "*")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: getEnvAsBool("ALLOW_CREDENTIALS", true),
	}

	return config
}

// getEnv 從環境變數獲取字串，如果不存在則使用默認值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 從環境變數獲取整數，如果不存在或無效則使用默認值
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool 從環境變數獲取布爾值，如果不存在或無效則使用默認值
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
