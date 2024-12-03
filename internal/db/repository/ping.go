package repository

import "context"

// PingRepository сервис для проверки состояния БД
type PingRepository interface {
	PingDB(ctx context.Context) error
}
