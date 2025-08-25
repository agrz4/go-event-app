package env

import (
	"os"
	"strconv"
)

// RedisConfig menyimpan konfigurasi Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// GetRedisConfig mengambil konfigurasi Redis dari environment variables
func GetRedisConfig() *RedisConfig {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbStr := getEnv("REDIS_DB", "0")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	return &RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}

// GetRedisAddr mendapatkan alamat Redis lengkap
func (r *RedisConfig) GetRedisAddr() string {
	return r.Host + ":" + r.Port
}

// getEnv helper function untuk mengambil environment variable dengan default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
