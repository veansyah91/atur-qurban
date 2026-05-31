package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/internal/model"
)

// MockContactCategoryRepository adalah mock untuk interface ContactCategoryRepository
type MockContactCategoryRepository struct {
	mock.Mock
}

func (m *MockContactCategoryRepository) Create(ctx context.Context, category *model.ContactCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockContactCategoryRepository) FindByTenantID(ctx context.Context, tenantID string) ([]*model.ContactCategory, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ContactCategory), args.Error(1)
}

func (m *MockContactCategoryRepository) FindByID(ctx context.Context, id string) (*model.ContactCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContactCategory), args.Error(1)
}

func (m *MockContactCategoryRepository) Update(ctx context.Context, category *model.ContactCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockContactCategoryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
