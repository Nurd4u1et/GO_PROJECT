package service

import (
	"context"
	"testing"

	"clinic-cli/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}
type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string, role model.Role) (*model.User, error)
}

func (m *MockUserRepo) Create(ctx context.Context, email, passwordHash string, role model.Role) (*model.User, error) {
	args := m.Called(ctx, email, passwordHash, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	return nil, nil
}

func TestAuthService_Register(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewAuthService(repo, "secret")

	repo.On(
		"Create",
		mock.Anything,
		"test@example.com",
		mock.AnythingOfType("string"),
		model.RolePatient,
	).Return(&model.User{
		ID:    1,
		Email: "test@example.com",
	}, nil)

	user, err := service.Register(
		context.Background(),
		"test@example.com",
		"password123",
		model.RolePatient,
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	repo.AssertExpectations(t)
}
