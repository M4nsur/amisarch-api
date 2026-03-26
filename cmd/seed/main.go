package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, читаем из окружения")
	}

	email := os.Getenv("SEED_EMAIL")
	password := os.Getenv("SEED_PASSWORD")
	dbURL := os.Getenv("DATABASE_URL")

	if email == "" || password == "" || dbURL == "" {
		log.Fatal("SEED_EMAIL, SEED_PASSWORD и DATABASE_URL обязательны")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal("подключение к БД:", err)
		defer conn.Close(ctx)
	}

	var exists bool
	err = conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email,
	).Scan(&exists)
	if err != nil {
		log.Fatal("ошибка проверки пользователя:", err)
	}

	if exists {
		log.Println("пользователь уже существует, пропускаем")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("ошибка хеширования пароля:", err)
	}

	_, err = conn.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, role, created_at)
         VALUES (gen_random_uuid(), $1, $2, 'owner', NOW())`,
		email, string(hash),
	)
	if err != nil {
		log.Fatal("ошибка создания пользователя:", err)
	}

	log.Printf("owner создан: %s\n", email)
}
