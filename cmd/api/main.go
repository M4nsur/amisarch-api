package main

import (
	"context"
	"net/http"
	"os"

	"github.com/M4nsur/amisarch-api/internal/handler"
	"github.com/M4nsur/amisarch-api/internal/repository"
	"github.com/M4nsur/amisarch-api/internal/service"
	"github.com/M4nsur/amisarch-api/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()

	isDev := os.Getenv("APP_ENV") != "production"
	if err := logger.Init(isDev); err != nil {
		panic(err)
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("не удалось подключиться к БД", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal("БД недоступна", zap.Error(err))
	}
	logger.Info("подключение к БД успешно")

	jwtSecret := os.Getenv("JWT_SECRET")
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtSecret)
	userHandler := handler.NewUsersHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", userHandler.Login)
	mux.HandleFunc("POST /auth/logout", userHandler.Logout)
	mux.HandleFunc("POST /auth/refresh", userHandler.Refresh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("сервер запущен", zap.String("port", port))
	http.ListenAndServe(":"+port, mux)
}
