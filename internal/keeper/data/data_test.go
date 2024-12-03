package data_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-keeper/internal/keeper/data"
	"go-keeper/internal/models"
)

// Мок для DataRepository
type MockDataRepository struct {
	mock.Mock
}

func (m *MockDataRepository) AddData(
	ctx context.Context, userID, dataType string, body []byte, meta string,
) (string, error) {
	args := m.Called(ctx, userID, dataType, body, meta)
	return args.String(0), args.Error(1)
}

func (m *MockDataRepository) GetAllDataByUserID(
	ctx context.Context, userID, dataType string,
) ([]*models.Data, error) {
	args := m.Called(ctx, userID, dataType)
	return args.Get(0).([]*models.Data), args.Error(1)
}

func (m *MockDataRepository) GetByUserID(
	ctx context.Context, userID, dataID, dataType string,
) (*models.Data, error) {
	args := m.Called(ctx, userID, dataID, dataType)
	return args.Get(0).(*models.Data), args.Error(1)
}

func (m *MockDataRepository) DeleteByID(
	ctx context.Context, userID, dataID, dataType string,
) error {
	args := m.Called(ctx, userID, dataID, dataType)
	return args.Error(0)
}

func (m *MockDataRepository) UpdateData(
	ctx context.Context, userID, dataID, dataType string, encryptedData []byte, meta string,
) error {
	args := m.Called(ctx, userID, dataID, dataType, encryptedData, meta)
	return args.Error(0)
}

type MockCrypt struct {
	mock.Mock
}

func (m *MockCrypt) Encrypt(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCrypt) Decrypt(data []byte) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func TestService_Add(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	dataType := data.LogPassDataType
	meta := "some meta"
	body := models.LogPassData{
		Login:    "user",
		Password: "password",
	}
	mockRepo := new(MockDataRepository)
	mockCrypt := new(MockCrypt)

	service := data.NewService(mockRepo, mockCrypt)

	t.Run(
		"Success", func(t *testing.T) {
			jsonData, _ := json.Marshal(body)
			encryptedData := []byte("encrypted data")
			expectedID := "data-id"

			mockCrypt.On("Encrypt", jsonData).Return(encryptedData, nil).Once()
			mockRepo.On("AddData", ctx, userID, dataType, encryptedData, meta).Return(
				expectedID, nil,
			).Once()

			id, err := service.Add(ctx, userID, dataType, body, meta)

			assert.NoError(t, err)
			assert.Equal(t, expectedID, id)
			mockCrypt.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		},
	)

	t.Run(
		"Error in Marshal", func(t *testing.T) {
			invalidBody := make(chan int) // Неподдерживаемый для маршализации тип

			_, err := service.Add(ctx, userID, dataType, invalidBody, meta)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "json: unsupported type")
			mockCrypt.AssertNotCalled(t, "Encrypt")
			mockRepo.AssertNotCalled(t, "AddData")
		},
	)

	t.Run(
		"Error in AddData", func(t *testing.T) {
			jsonData, _ := json.Marshal(body)
			encryptedData := []byte("encrypted data")
			addDataErr := errors.New("database error")

			mockCrypt.On("Encrypt", jsonData).Return(encryptedData, nil).Once()
			mockRepo.On("AddData", ctx, userID, dataType, encryptedData, meta).Return(
				"", addDataErr,
			).Once()

			_, err := service.Add(ctx, userID, dataType, body, meta)

			assert.Error(t, err)
			assert.Equal(t, addDataErr, err)
			mockCrypt.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		},
	)
}

func TestService_GetAllData(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	dataType := data.LogPassDataType

	mockRepo := new(MockDataRepository)
	mockCrypt := new(MockCrypt)

	service := data.NewService(mockRepo, mockCrypt)

	t.Run(
		"Success", func(t *testing.T) {
			mockData := []*models.Data{
				{
					Data: "encrypted data 1",
					Meta: "meta1",
					Type: data.LogPassDataType,
					ID:   "1",
				},
				{
					Data: "encrypted data 2",
					Meta: "meta2",
					Type: data.LogPassDataType,
					ID:   "2",
				},
			}

			decryptedData1 := `{"login": "user1", "password": "pass1"}`
			decryptedData2 := `{"login": "user2", "password": "pass2"}`

			parsedData1 := models.LogPassData{
				Login:    "user1",
				Password: "pass1",
			}
			parsedData2 := models.LogPassData{
				Login:    "user2",
				Password: "pass2",
			}

			mockRepo.On("GetAllDataByUserID", ctx, userID, dataType).Return(mockData, nil).Once()
			mockCrypt.On("Decrypt", []byte("encrypted data 1")).Return(decryptedData1, nil).Once()
			mockCrypt.On("Decrypt", []byte("encrypted data 2")).Return(decryptedData2, nil).Once()

			result, err := service.GetAllData(ctx, userID, dataType)

			assert.NoError(t, err)
			assert.Len(t, result, 2)

			assert.Equal(t, parsedData1, result[0].Data)
			assert.Equal(t, "meta1", result[0].Meta)
			assert.Equal(t, data.LogPassDataType, result[0].Type)
			assert.Equal(t, "1", result[0].ID)

			assert.Equal(t, parsedData2, result[1].Data)
			assert.Equal(t, "meta2", result[1].Meta)
			assert.Equal(t, data.LogPassDataType, result[1].Type)
			assert.Equal(t, "2", result[1].ID)

			mockRepo.AssertExpectations(t)
			mockCrypt.AssertExpectations(t)
		},
	)

}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	dataID := "data-id"
	dataType := data.LogPassDataType

	mockRepo := new(MockDataRepository)
	mockCrypt := new(MockCrypt)

	service := data.NewService(mockRepo, mockCrypt)

	t.Run(
		"Success", func(t *testing.T) {
			mockData := &models.Data{
				Data: "encrypted data",
				Meta: "meta",
				Type: data.LogPassDataType,
				ID:   dataID,
			}

			decryptedData := `{"login": "user1", "password": "pass1"}`
			parsedData := models.LogPassData{
				Login:    "user1",
				Password: "pass1",
			}

			mockRepo.On("GetByUserID", ctx, userID, dataID, dataType).Return(mockData, nil).Once()
			mockCrypt.On("Decrypt", []byte("encrypted data")).Return(decryptedData, nil).Once()

			result, err := service.GetByID(ctx, userID, dataID, dataType)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, parsedData, result.Data)
			assert.Equal(t, "meta", result.Meta)
			assert.Equal(t, data.LogPassDataType, result.Type)
			assert.Equal(t, dataID, result.ID)

			mockRepo.AssertExpectations(t)
			mockCrypt.AssertExpectations(t)
		},
	)

}

func TestService_DeleteByID(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	dataID := "data-id"
	dataType := data.LogPassDataType

	mockRepo := new(MockDataRepository)
	service := data.NewService(mockRepo, nil)

	t.Run(
		"Success", func(t *testing.T) {
			mockRepo.On("DeleteByID", ctx, userID, dataID, dataType).Return(nil).Once()

			err := service.DeleteByID(ctx, userID, dataID, dataType)

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		},
	)

	t.Run(
		"Error in DeleteByID", func(t *testing.T) {
			mockRepo.On(
				"DeleteByID", ctx, userID, dataID, dataType,
			).Return(errors.New("delete error")).Once()

			err := service.DeleteByID(ctx, userID, dataID, dataType)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "delete error")
			mockRepo.AssertExpectations(t)
		},
	)
}

func TestParseJsonData(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		dataType  string
		expected  interface{}
		expectErr bool
	}{
		{
			name:     "Valid LogPassData",
			data:     []byte(`{"login":"user","password":"pass"}`),
			dataType: data.LogPassDataType,
			expected: models.LogPassData{
				Login:    "user",
				Password: "pass",
			},
			expectErr: false,
		},
		{
			name:     "Valid CardData",
			data:     []byte(`{"number":"1234567812345678","expired_at":"12/24","cvv":"123"}`),
			dataType: data.CardDataType,
			expected: models.CardData{
				Number:    "1234567812345678",
				ExpiredAt: "12/24",
				CVV:       "123",
			},
			expectErr: false,
		},
		{
			name:     "Valid TextData",
			data:     []byte(`{"text":"Hello, World!"}`),
			dataType: data.TextDataType,
			expected: models.TextData{
				Text: "Hello, World!",
			},
			expectErr: false,
		},
		{
			name:      "Invalid JSON for LogPassData",
			data:      []byte(`{"login":"user","password":123}`),
			dataType:  data.LogPassDataType,
			expectErr: true,
		},
		{
			name:      "Unknown DataType",
			data:      []byte(`{"some_field":"some_value"}`),
			dataType:  "UnknownType",
			expectErr: true,
		},
		{
			name:      "Invalid JSON Syntax",
			data:      []byte(`{"login":"user",}`),
			dataType:  data.LogPassDataType,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result, err := data.ParseJsonData(tt.data, tt.dataType)

				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, result)
				}
			},
		)
	}
}
