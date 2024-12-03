package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"go-keeper/internal/db/repository"
	"go-keeper/internal/http/rest/schemas"
	"go-keeper/internal/models"
	"go-keeper/pkg/hashing"
)

const (
	LogPassDataType = "logpass"
	CardDataType    = "card"
	TextDataType    = "text"
	FileDataType    = "file"
)

var ErrTypeNotFound = fmt.Errorf("data type not found")

// Service Сервис для работы с данными.
type Service struct {
	dataRepo repository.DataRepository
	crypt    hashing.CryptInterface
}

// NewService создает новый экземпляр Service.
func NewService(dataRepo repository.DataRepository, crypt hashing.CryptInterface) *Service {
	return &Service{dataRepo: dataRepo, crypt: crypt}
}

// Add добавить новые данные
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataType: Тип данных.
// body: Данные для добавления.
// meta: Метаданные.
//
// Возвращает ID добавленных данных и ошибку, если таковая возникла.
func (s *Service) Add(
	ctx context.Context, userID string, dataType string, body interface{}, meta string,
) (string, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		slog.Error("error to marshall data", slog.String("err", err.Error()))
		return "", err
	}

	encrypted, err := s.crypt.Encrypt(jsonData)
	if err != nil {
		slog.Error("error to encrypt data", slog.String("err", err.Error()))
		return "", err
	}
	id, err := s.dataRepo.AddData(ctx, userID, dataType, encrypted, meta)
	if err != nil {
		slog.Error("error to add data to db", slog.String("err", err.Error()))
		return "", err
	}

	return id, nil
}

// GetAllData получает все данные пользователя определенного типа.
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataType: Тип данных.
//
// Возвращает слайс структур schemas.DataResponse и ошибку, если таковая возникла.
func (s *Service) GetAllData(
	ctx context.Context, userID string, dataType string,
) ([]*schemas.DataResponse, error) {
	dataSet, err := s.dataRepo.GetAllDataByUserID(ctx, userID, dataType)
	if err != nil {
		return nil, err
	}

	var dataResponses []*schemas.DataResponse
	for _, data := range dataSet {
		decrypted, err := s.crypt.Decrypt([]byte(data.Data))
		if err != nil {
			return nil, err
		}
		parsedData, err := ParseJsonData([]byte(decrypted), data.Type)
		if err != nil {
			return nil, err
		}

		dataResponses = append(
			dataResponses, &schemas.DataResponse{
				Data: parsedData,
				Meta: data.Meta,
				Type: data.Type,
				ID:   data.ID,
			},
		)
	}

	return dataResponses, nil
}

// GetByID получает данные по ID.
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataID: ID данных.
// dataType: Тип данных.
//
// Возвращает указатель на структуру schemas.DataResponse и ошибку, если таковая возникла.
func (s *Service) GetByID(
	ctx context.Context, userID string, dataID string, dataType string,
) (*schemas.DataResponse, error) {
	data, err := s.dataRepo.GetByUserID(ctx, userID, dataID, dataType)
	if err != nil {
		return nil, err
	}
	decrypted, err := s.crypt.Decrypt([]byte(data.Data))
	if err != nil {
		return nil, err
	}
	parsedData, err := ParseJsonData([]byte(decrypted), data.Type)
	if err != nil {
		return nil, err
	}

	return &schemas.DataResponse{
		Data: parsedData,
		Meta: data.Meta,
		Type: data.Type,
		ID:   data.ID,
	}, nil
}

// DeleteByID удаляет данные по ID.
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataID: ID данных.
// dataType: Тип данных.
//
// Возвращает ошибку, если таковая возникла.
func (s *Service) DeleteByID(
	ctx context.Context, userID string, dataID string, dataType string,
) error {
	return s.dataRepo.DeleteByID(ctx, userID, dataID, dataType)
}

// Update обновляет данные.
//
// ctx: Контекст запроса.
// userID: ID пользователя.
// dataID: ID данных.
// dataType: Тип данных.
// data: Данные для обновления.
// meta: Метаданные.
//
// Возвращает ошибку, если таковая возникла.
func (s *Service) Update(
	ctx context.Context,
	userID string,
	dataID string,
	dataType string,
	data interface{},
	meta string,
) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("error to marshall data", slog.String("err", err.Error()))
		return err
	}

	encrypted, err := s.crypt.Encrypt(jsonData)
	if err != nil {
		slog.Error("error to encrypt data", slog.String("err", err.Error()))
		return err
	}

	err = s.dataRepo.UpdateData(ctx, userID, dataID, dataType, encrypted, meta)
	if err != nil {
		slog.Error("error to update data to db", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// ParseJsonData разбирает JSON данные в зависимости от типа данных.
//
// data: JSON данные.
// dataType: Тип данных.
//
// Возвращает интерфейс, содержащий разобранные данные, и ошибку, если таковая возникла.
func ParseJsonData(data []byte, dataType string) (interface{}, error) {
	switch dataType {
	case LogPassDataType:
		var parsedData models.LogPassData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	case CardDataType:
		var parsedData models.CardData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	case TextDataType:
		var parsedData models.TextData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	case FileDataType:
		var parsedData models.FileData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	}
	return nil, ErrTypeNotFound
}
