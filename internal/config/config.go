package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	AppEnv       string
	Port         string
	SeedEmail    string
	SeedPassword string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, читаем из окружения")
	}

	return &Config{
		DatabaseURL:  mustGet("DATABASE_URL"),
		JWTSecret:    mustGet("JWT_SECRET"),
		AppEnv:       getOrDefault("APP_ENV", "development"),
		Port:         getOrDefault("PORT", "8080"),
		SeedEmail:    mustGet("SEED_EMAIL"),
		SeedPassword: mustGet("SEED_PASSWORD"),
	}
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("обязательная переменная %s не задана", key)
	}
	return v
}

func getOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
