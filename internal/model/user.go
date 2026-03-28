package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"            db:"id"`
	Email        string    `json:"email"         db:"email"`
	PasswordHash string    `json:"-"             db:"password_hash"`
	Role         string    `json:"role"          db:"role"`
	CreatedAt    time.Time `json:"created_at"    db:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
