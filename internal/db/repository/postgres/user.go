package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"go-keeper/internal/db/repository/errs"
	"go-keeper/internal/models"
)

// UserRepository реализует репозиторий пользователей для PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository.
//
// db: Экземпляр базы данных.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя в базе данных.
//
// ctx: Контекст запроса.
// user: Указатель на структуру models.User, содержащую данные нового пользователя.
//
// Возвращает указатель на ID созданного пользователя и ошибку, если таковая возникла.
func (u *UserRepository) Create(ctx context.Context, user *models.User) (
	*string, error,
) {
	var userID string

	query := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id
	`
	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Login,
		user.Password,
	).Scan(&userID)

	if err != nil {
		slog.Error("error to insert user", slog.String("err", err.Error()))
		return nil, err
	}

	return &userID, nil
}

// GetUserByLogin получает пользователя из базы данных по его логину.
//
// ctx: Контекст запроса.
// login: Логин пользователя.
//
// Возвращает указатель на структуру models.User и ошибку, если таковая возникла.
// Возвращает errs.ErrUserNotFound, если пользователь не найден.
func (u *UserRepository) GetUserByLogin(ctx context.Context, login string) (
	*models.User, error,
) {
	var user models.User

	query := `
		SELECT id, login, password
		FROM users
		WHERE login = $1
	`

	row := u.db.QueryRowContext(
		ctx,
		query,
		login,
	)

	if err := row.Scan(&user.ID, &user.Login, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		slog.Error("error to scan user from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &user, nil
}
