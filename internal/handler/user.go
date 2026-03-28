package handler

import (
	"context"

	"github.com/M4nsur/amisarch-api/internal/model"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService UserService
}

type UserService interface {
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

func NewUsersHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
