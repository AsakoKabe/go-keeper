package repository

import (
	"context"

	"go-keeper/internal/models"
)

// UserRepository определяет интерфейс для работы с репозиторием пользователей.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*string, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}
