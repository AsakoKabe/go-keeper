package repository

import (
	"context"

	"go-keeper/internal/models"
)

// DataRepository репозиторий для работы с данными
type DataRepository interface {
	AddData(
		ctx context.Context, userID string, dataType string, body []byte, meta string,
	) (string, error)

	GetAllDataByUserID(ctx context.Context, userID string, dataType string) ([]*models.Data, error)
	GetByUserID(ctx context.Context, userID string, dataID string, dataType string) (
		*models.Data, error,
	)
	DeleteByID(ctx context.Context, userID string, dataID string, dataType string) error
	UpdateData(
		ctx context.Context,
		userID string,
		dataID string,
		dataType string,
		encryptedData []byte,
		meta string,
	) error
}
