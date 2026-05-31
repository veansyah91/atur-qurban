package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/internal/model"
)

// MockContactRepository adalah mock untuk interface ContactRepository
type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) Create(ctx context.Context, contact *model.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) FindByTenantID(ctx context.Context, tenantID string) ([]*model.Contact, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Contact), args.Error(1)
}

func (m *MockContactRepository) FindByID(ctx context.Context, id string) (*model.Contact, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Contact), args.Error(1)
}

func (m *MockContactRepository) Update(ctx context.Context, contact *model.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
