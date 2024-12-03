package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"go-keeper/internal/db/repository/errs"
	"go-keeper/internal/models"
)

// DataRepository структура реализации сервиса к БД для postgres
type DataRepository struct {
	db *sql.DB
}

// NewDataRepository конструктор для DataRepository
func NewDataRepository(db *sql.DB) *DataRepository {
	return &DataRepository{db: db}
}

// AddData добавляет новые данные в базу данных.
//
// ctx: Контекст запроса.
// userID: ID пользователя, которому принадлежат данные.
// dataType: Тип данных.
// data: Данные в байтовом представлении.
// meta: Метаданные.
//
// Возвращает ID добавленной записи и ошибку, если таковая возникла.
func (d *DataRepository) AddData(
	ctx context.Context, userID string, dataType string, data []byte, meta string,
) (string, error) {
	query := `
		INSERT INTO user_data (user_id, data_type, data, meta)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var insertedID string

	err := d.db.QueryRowContext(
		ctx, query, userID, dataType, data, meta,
	).Scan(&insertedID)
	if err != nil {
		return "", err
	}

	return insertedID, nil
}

// GetAllDataByUserID получает все данные по типу данных и ID пользователя
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataType: Тип данных.
//
// Возвращает слайс структур models.Data и ошибку, если таковая возникла.
func (d *DataRepository) GetAllDataByUserID(
	ctx context.Context, userID string, dataType string,
) ([]*models.Data, error) {
	query := `
		SELECT id, data_type, data, meta FROM user_data
		WHERE user_id = $1 and data_type = $2 and is_deleted = false 
	`

	rows, err := d.db.QueryContext(ctx, query, userID, dataType)
	if err != nil {
		return nil, err
	}

	var dataSet []*models.Data
	for rows.Next() {
		var data models.Data
		if err = rows.Scan(&data.ID, &data.Type, &data.Data, &data.Meta); err != nil {
			return nil, err
		}

		dataSet = append(dataSet, &data)
	}

	return dataSet, nil
}

// GetByUserID получает данные пользователя по его ID, ID данных и типу данных.
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataID: ID данных.
// dataType: Тип данных.
//
// Возвращает указатель на структуру models.Data и ошибку, если таковая возникла.
// Возвращает errs.ErrDataNotFound, если данные не найдены.
func (d *DataRepository) GetByUserID(
	ctx context.Context, userID string, dataID string, dataType string,
) (*models.Data, error) {
	query := `
		SELECT id, data_type, data, meta FROM user_data
		WHERE user_id = $1 and id = $2 and data_type = $3 and is_deleted = false
	`

	var data models.Data

	row := d.db.QueryRowContext(ctx, query, userID, dataID, dataType)
	if err := row.Scan(&data.ID, &data.Type, &data.Data, &data.Meta); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		slog.Error("error to scan data from db", slog.String("err", err.Error()))
		return nil, err
	}

	return &data, nil

}

// DeleteByID помечает данные как удаленные по ID пользователя, ID данных и типу данных.
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataID: ID данных.
// dataType: Тип данных.
//
// Возвращает ошибку, если таковая возникла.
// Возвращает errs.ErrDataNotFound, если данные не найдены.
func (d *DataRepository) DeleteByID(
	ctx context.Context, userID string, dataID string, dataType string,
) error {
	query := `
		UPDATE user_data SET is_deleted = true
		WHERE user_id = $1 and id = $2 and data_type = $3 and is_deleted = false
	`

	result, err := d.db.ExecContext(ctx, query, userID, dataID, dataType)
	if err != nil {
		return fmt.Errorf("unable to set deleted data: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrDataNotFound
	}

	return nil
}

// UpdateData обновляет данные в базе данных.
//
// ctx: Контекст запроса.
// userID: ID пользователя, которому принадлежат данные.
// dataID: ID данных для обновления.
// dataType: Тип данных.
// encryptedData: Обновленные данные в байтовом представлении.
// meta: Обновленные метаданные.
//
// Возвращает ошибку, если таковая возникла.
// Возвращает errs.ErrDataNotFound, если данные не найдены.
func (d *DataRepository) UpdateData(
	ctx context.Context,
	userID string,
	dataID string,
	dataType string,
	encryptedData []byte,
	meta string,
) error {
	query := `
		UPDATE user_data SET data = $1, meta = $2 WHERE id = $3 and data_type = $4 and user_id = $5 and is_deleted = false
	`

	result, err := d.db.ExecContext(
		ctx, query, encryptedData, meta, dataID, dataType, userID,
	)
	if err != nil {
		return fmt.Errorf("unable to update data: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrDataNotFound
	}

	return nil
}
