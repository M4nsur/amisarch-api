package main

import (
	"context"
	"log"

	"github.com/M4nsur/amisarch-api/internal/config"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("подключение к БД:", err)
	}
	defer conn.Close(context.Background())

	var exists bool
	conn.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", cfg.SeedEmail,
	).Scan(&exists)

	if exists {
		log.Println("пользователь уже существует, пропускаем")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("ошибка хеширования пароля:", err)
	}

	_, err = conn.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, role, created_at)
		 VALUES (gen_random_uuid(), $1, $2, 'owner', NOW())`,
		cfg.SeedEmail, string(hash),
	)
	if err != nil {
		log.Fatal("ошибка создания пользователя:", err)
	}

	log.Printf("owner создан: %s\n", cfg.SeedEmail)
}
