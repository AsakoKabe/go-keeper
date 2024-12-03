package user

import (
	"context"
	"testing"

	"go-keeper/internal/db/repository/errs"
	"go-keeper/internal/models"
	"go-keeper/pkg/hashing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByLogin(ctx context.Context, login string) (
	*models.User, error,
) {
	args := m.Called(ctx, login)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*string, error) {
	args := m.Called(ctx, user)
	if id, ok := args.Get(0).(*string); ok {
		return id, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	user := &models.User{
		Login:    "testuser",
		Password: "password123",
	}
	expectedID := "12345"

	mockRepo.On("GetUserByLogin", ctx, user.Login).Return(nil, errs.ErrUserNotFound)
	mockRepo.On("Create", ctx, user).Return(&expectedID, nil)

	id, err := service.CreateUser(ctx, user)

	require.NoError(t, err)
	require.Equal(t, expectedID, *id)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_ExistingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	user := &models.User{
		Login:    "testuser",
		Password: "password123",
	}

	mockRepo.On("GetUserByLogin", ctx, user.Login).Return(&models.User{Login: user.Login}, nil)

	id, err := service.CreateUser(ctx, user)

	require.ErrorIs(t, err, errs.ErrLoginAlreadyExist)
	require.Nil(t, id)

	mockRepo.AssertExpectations(t)
}

func TestCorrectCredentials(t *testing.T) {
	service := NewService(nil)

	user := &models.User{Password: "password123"}
	hashedPassword, _ := hashing.HashPassword("password123")
	existingUser := &models.User{Password: hashedPassword}

	require.True(t, service.CorrectCredentials(user, existingUser))

	user.Password = "wrongpassword"
	require.False(t, service.CorrectCredentials(user, existingUser))
}
