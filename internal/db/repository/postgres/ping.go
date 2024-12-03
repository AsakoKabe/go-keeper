package postgres

import (
	"context"
	"database/sql"
)

// PingRepository структура реализации сервиса к БД для postgres
type PingRepository struct {
	db *sql.DB
}

// NewPingRepository конструктор для PingRepository
func NewPingRepository(db *sql.DB) *PingRepository {
	return &PingRepository{db: db}
}

// PingDB функция для отправки пинга в БД для postgres
func (p *PingRepository) PingDB(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
