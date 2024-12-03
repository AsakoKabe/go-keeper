package user

import (
	"context"
	"errors"
	"log/slog"

	"go-keeper/internal/db/repository"
	"go-keeper/internal/db/repository/errs"
	"go-keeper/internal/models"
	"go-keeper/pkg/hashing"
)

// Service Сервис работы с пользователями.
type Service struct {
	userRepo repository.UserRepository
}

// NewService создает новый экземпляр Service.
//
// userRepo: Репозиторий пользователей.
func NewService(userRepo repository.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

// CreateUser создает нового пользователя.
//
// ctx: Контекст запроса.
// user: Указатель на структуру models.User, содержащую данные нового пользователя.
//
// Возвращает указатель на ID созданного пользователя и ошибку, если таковая возникла.
// Возвращает errs.ErrLoginAlreadyExist, если пользователь с таким логином уже существует.
func (s *Service) CreateUser(ctx context.Context, user *models.User) (*string, error) {
	exitedUser, err := s.userRepo.GetUserByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		slog.Error("error to check user exist", slog.String("err", err.Error()))
		return nil, err
	}
	if exitedUser != nil {
		slog.Error("user already exist")
		return nil, errs.ErrLoginAlreadyExist
	}

	user.Password, err = hashing.HashPassword(user.Password)
	if err != nil {
		slog.Error("error to hashing password", slog.String("err", err.Error()))
		return nil, err
	}
	return s.userRepo.Create(ctx, user)
}

// GetUser получает пользователя по логину.
//
// ctx: Контекст запроса.
// user: Указатель на структуру models.User, содержащую логин пользователя.
//
// Возвращает указатель на структуру models.User и ошибку, если таковая возникла.
func (s *Service) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	return s.userRepo.GetUserByLogin(ctx, user.Login)
}

// CorrectCredentials проверяет корректность учетных данных пользователя.
//
// user: Указатель на структуру models.User, содержащую введенные пользователем учетные данные.
// existedUser: Указатель на структуру models.User, содержащую данные пользователя из базы данных.
//
// Возвращает true, если учетные данные корректны, иначе false.
func (s *Service) CorrectCredentials(user *models.User, existedUser *models.User) bool {
	return hashing.VerifyPassword(user.Password, existedUser.Password)
}
