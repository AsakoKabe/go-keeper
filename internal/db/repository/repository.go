package repository

import (
	"database/sql"

	"go-keeper/internal/db/repository/postgres"
)

// Repositories Набор репозиториев для доступов к БД
type Repositories struct {
	PingRepository PingRepository
	UserRepository UserRepository
	DataRepository DataRepository
}

// NewPostgresRepositories функция для создания сервисов postgres
func NewPostgresRepositories(db *sql.DB) (*Repositories, error) {
	return &Repositories{
		PingRepository: postgres.NewPingRepository(db),
		UserRepository: postgres.NewUserRepository(db),
		DataRepository: postgres.NewDataRepository(db),
	}, nil
}
