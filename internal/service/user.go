package service

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/M4nsur/amisarch-api/internal/model"
	"github.com/M4nsur/amisarch-api/internal/repository"
	"github.com/M4nsur/amisarch-api/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("Login: ошибка поиска пользователя", zap.String("email", req.Email), zap.Error(err))
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}
	if user == nil {
		return nil, errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	accessToken, err := s.generateToken(user, 15*time.Minute, "access")
	if err != nil {
		logger.Error("Login: ошибка генерации access token", zap.Error(err))
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	refreshToken, err := s.generateToken(user, 30*24*time.Hour, "refresh")
	if err != nil {
		logger.Error("Login: ошибка генерации refresh token", zap.Error(err))
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	logger.Info("Login: успешный вход", zap.String("email", user.Email), zap.String("role", user.Role))

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {

	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, errors.New("невалидный токен")
	}

	if claims["type"] != "refresh" {
		return nil, errors.New("невалидный токен")
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, errors.New("невалидный токен")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("пользователь не найден")
	}

	accessToken, err := s.generateToken(user, 15*time.Minute, "access")
	if err != nil {
		logger.Error("Refresh: ошибка генерации access token", zap.Error(err))
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	newRefreshToken, err := s.generateToken(user, 30*24*time.Hour, "refresh")
	if err != nil {
		logger.Error("Refresh: ошибка генерации refresh token", zap.Error(err))
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Error("GetByID: ошибка", zap.String("id", id.String()), zap.Error(err))
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

func (s *userService) generateToken(user *model.User, duration time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": user.Role,
		"type": tokenType,
		"exp":  time.Now().Add(duration).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *userService) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("невалидный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("невалидный токен")
	}

	return claims, nil
}
