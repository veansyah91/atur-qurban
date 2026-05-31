package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/internal/model"
)

// MockContactCategoryRepository mock untuk ContactCategoryRepository
type MockContactCategoryRepository struct {
	mock.Mock
}

func (m *MockContactCategoryRepository) Create(ctx context.Context, category *model.ContactCategory) error {
	return m.Called(ctx, category).Error(0)
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
	return m.Called(ctx, category).Error(0)
}

func (m *MockContactCategoryRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockContactRepository mock untuk ContactRepository
type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) Create(ctx context.Context, contact *model.Contact) error {
	return m.Called(ctx, contact).Error(0)
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
	return m.Called(ctx, contact).Error(0)
}

func (m *MockContactRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ─────────────────────────────────────────────
// Test ContactCategoryService
// ─────────────────────────────────────────────

func TestContactCategoryService_Create(t *testing.T) {
	tests := []struct {
		name     string
		tenantID string
		catName  string
		mock     func(*MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			tenantID: "tenant-1",
			catName:  "Hewan Sapi",
			mock: func(m *MockContactCategoryRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*model.ContactCategory")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "tenantID kosong",
			tenantID: "",
			catName:  "Hewan Sapi",
			mock:     func(m *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "name kosong",
			tenantID: "tenant-1",
			catName:  "",
			mock:     func(m *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "repo error",
			tenantID: "tenant-1",
			catName:  "Hewan Sapi",
			mock: func(m *MockContactCategoryRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*model.ContactCategory")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactCategoryRepository)
			tt.mock(repo)
			svc := NewContactCategoryService(repo)
			result, err := svc.Create(context.Background(), tt.tenantID, tt.catName)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
				assert.Equal(t, tt.tenantID, result.TenantID)
				assert.Equal(t, tt.catName, result.Name)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestContactCategoryService_GetByTenant(t *testing.T) {
	tests := []struct {
		name     string
		tenantID string
		mock     func(*MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByTenantID", mock.Anything, "tenant-1").Return([]*model.ContactCategory{
					{ID: "cat-1", TenantID: "tenant-1", Name: "Umum"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "tenantID kosong",
			tenantID: "",
			mock:     func(m *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "repo error",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByTenantID", mock.Anything, "tenant-1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactCategoryRepository)
			tt.mock(repo)
			svc := NewContactCategoryService(repo)
			result, err := svc.GetByTenant(context.Background(), tt.tenantID)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestContactCategoryService_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		tenantID string
		mock     func(*MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			id:       "cat-1",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1", Name: "Umum"}, nil,
				)
			},
			wantErr: false,
		},
		{
			name:     "id kosong",
			id:       "",
			tenantID: "tenant-1",
			mock:     func(m *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "tidak ditemukan",
			id:       "cat-x",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "repo error",
			id:       "cat-1",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactCategoryRepository)
			tt.mock(repo)
			svc := NewContactCategoryService(repo)
			result, err := svc.GetByID(context.Background(), tt.id, tt.tenantID)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestContactCategoryService_Update(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		tenantID string
		catName  string
		mock     func(*MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			id:       "cat-1",
			tenantID: "tenant-1",
			catName:  "Nama Baru",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1", Name: "Lama"}, nil,
				)
				m.On("Update", mock.Anything, mock.MatchedBy(func(c *model.ContactCategory) bool {
					return c.Name == "Nama Baru"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "id kosong",
			id:       "",
			tenantID: "tenant-1",
			catName:  "Nama Baru",
			mock:     func(m *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "tidak ditemukan",
			id:       "cat-x",
			tenantID: "tenant-1",
			catName:  "Nama Baru",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "repo error saat update",
			id:       "cat-1",
			tenantID: "tenant-1",
			catName:  "Nama Baru",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1", Name: "Lama"}, nil,
				)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactCategoryRepository)
			tt.mock(repo)
			svc := NewContactCategoryService(repo)
			result, err := svc.Update(context.Background(), tt.id, tt.tenantID, tt.catName)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
				assert.Equal(t, tt.catName, result.Name)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestContactCategoryService_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		tenantID string
		mock     func(*MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			id:       "cat-1",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1"}, nil,
				)
				m.On("Delete", mock.Anything, "cat-1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "id kosong",
			id:       "",
			tenantID: "tenant-1",
			mock:     func(m *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "tidak ditemukan",
			id:       "cat-x",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "repo error",
			id:       "cat-1",
			tenantID: "tenant-1",
			mock: func(m *MockContactCategoryRepository) {
				m.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1"}, nil,
				)
				m.On("Delete", mock.Anything, "cat-1").Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactCategoryRepository)
			tt.mock(repo)
			svc := NewContactCategoryService(repo)
			err := svc.Delete(context.Background(), tt.id, tt.tenantID)
			assert.Equal(t, tt.wantErr, err != nil)
			repo.AssertExpectations(t)
		})
	}
}

// ─────────────────────────────────────────────
// Test ContactService
// ─────────────────────────────────────────────

func TestContactService_Create(t *testing.T) {
	tests := []struct {
		name       string
		tenantID   string
		categoryID string
		contName   string
		mock       func(*MockContactRepository, *MockContactCategoryRepository)
		wantErr    bool
	}{
		{
			name:       "happy path",
			tenantID:   "tenant-1",
			categoryID: "cat-1",
			contName:   "Ahmad Syukri",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				cr.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1"}, nil,
				)
				r.On("Create", mock.Anything, mock.AnythingOfType("*model.Contact")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "tenantID kosong",
			tenantID:   "",
			categoryID: "cat-1",
			contName:   "Ahmad Syukri",
			mock:       func(r *MockContactRepository, cr *MockContactCategoryRepository) {},
			wantErr:    true,
		},
		{
			name:       "name kosong",
			tenantID:   "tenant-1",
			categoryID: "cat-1",
			contName:   "",
			mock:       func(r *MockContactRepository, cr *MockContactCategoryRepository) {},
			wantErr:    true,
		},
		{
			name:       "kategori tidak ditemukan",
			tenantID:   "tenant-1",
			categoryID: "cat-x",
			contName:   "Ahmad Syukri",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				cr.On("FindByID", mock.Anything, "cat-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:       "repo error",
			tenantID:   "tenant-1",
			categoryID: "cat-1",
			contName:   "Ahmad Syukri",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				cr.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1"}, nil,
				)
				r.On("Create", mock.Anything, mock.AnythingOfType("*model.Contact")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactRepository)
			catRepo := new(MockContactCategoryRepository)
			tt.mock(repo, catRepo)
			svc := NewContactService(repo, catRepo)
			result, err := svc.Create(context.Background(), tt.tenantID, tt.categoryID, tt.contName, nil, nil)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
				assert.Equal(t, tt.contName, result.Name)
				assert.True(t, result.IsActive)
			}
			repo.AssertExpectations(t)
			catRepo.AssertExpectations(t)
		})
	}
}

func TestContactService_GetByTenant(t *testing.T) {
	tests := []struct {
		name     string
		tenantID string
		mock     func(*MockContactRepository, *MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByTenantID", mock.Anything, "tenant-1").Return([]*model.Contact{
					{ID: "con-1", TenantID: "tenant-1", Name: "Ahmad"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "tenantID kosong",
			tenantID: "",
			mock:     func(r *MockContactRepository, cr *MockContactCategoryRepository) {},
			wantErr:  true,
		},
		{
			name:     "repo error",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByTenantID", mock.Anything, "tenant-1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactRepository)
			catRepo := new(MockContactCategoryRepository)
			tt.mock(repo, catRepo)
			svc := NewContactService(repo, catRepo)
			result, err := svc.GetByTenant(context.Background(), tt.tenantID)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestContactService_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		tenantID string
		mock     func(*MockContactRepository, *MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			id:       "con-1",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(
					&model.Contact{ID: "con-1", TenantID: "tenant-1", Name: "Ahmad"}, nil,
				)
			},
			wantErr: false,
		},
		{
			name:     "tidak ditemukan",
			id:       "con-x",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "repo error",
			id:       "con-1",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactRepository)
			catRepo := new(MockContactCategoryRepository)
			tt.mock(repo, catRepo)
			svc := NewContactService(repo, catRepo)
			result, err := svc.GetByID(context.Background(), tt.id, tt.tenantID)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestContactService_Update(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		tenantID   string
		categoryID string
		contName   string
		mock       func(*MockContactRepository, *MockContactCategoryRepository)
		wantErr    bool
	}{
		{
			name:       "happy path",
			id:         "con-1",
			tenantID:   "tenant-1",
			categoryID: "cat-1",
			contName:   "Nama Baru",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(
					&model.Contact{ID: "con-1", TenantID: "tenant-1", Name: "Lama"}, nil,
				)
				cr.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1"}, nil,
				)
				r.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Contact) bool {
					return c.Name == "Nama Baru"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "contact tidak ditemukan",
			id:         "con-x",
			tenantID:   "tenant-1",
			categoryID: "cat-1",
			contName:   "Nama Baru",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:       "kategori tidak ditemukan",
			id:         "con-1",
			tenantID:   "tenant-1",
			categoryID: "cat-x",
			contName:   "Nama Baru",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(
					&model.Contact{ID: "con-1", TenantID: "tenant-1", Name: "Lama"}, nil,
				)
				cr.On("FindByID", mock.Anything, "cat-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:       "repo error saat update",
			id:         "con-1",
			tenantID:   "tenant-1",
			categoryID: "cat-1",
			contName:   "Nama Baru",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(
					&model.Contact{ID: "con-1", TenantID: "tenant-1", Name: "Lama"}, nil,
				)
				cr.On("FindByID", mock.Anything, "cat-1").Return(
					&model.ContactCategory{ID: "cat-1", TenantID: "tenant-1"}, nil,
				)
				r.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactRepository)
			catRepo := new(MockContactCategoryRepository)
			tt.mock(repo, catRepo)
			svc := NewContactService(repo, catRepo)
			result, err := svc.Update(context.Background(), tt.id, tt.tenantID, tt.categoryID, tt.contName, true, nil, nil)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, result)
				assert.Equal(t, tt.contName, result.Name)
			}
			repo.AssertExpectations(t)
			catRepo.AssertExpectations(t)
		})
	}
}

func TestContactService_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		tenantID string
		mock     func(*MockContactRepository, *MockContactCategoryRepository)
		wantErr  bool
	}{
		{
			name:     "happy path",
			id:       "con-1",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(
					&model.Contact{ID: "con-1", TenantID: "tenant-1"}, nil,
				)
				r.On("Delete", mock.Anything, "con-1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "contact tidak ditemukan",
			id:       "con-x",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-x").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "repo error",
			id:       "con-1",
			tenantID: "tenant-1",
			mock: func(r *MockContactRepository, cr *MockContactCategoryRepository) {
				r.On("FindByID", mock.Anything, "con-1").Return(
					&model.Contact{ID: "con-1", TenantID: "tenant-1"}, nil,
				)
				r.On("Delete", mock.Anything, "con-1").Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockContactRepository)
			catRepo := new(MockContactCategoryRepository)
			tt.mock(repo, catRepo)
			svc := NewContactService(repo, catRepo)
			err := svc.Delete(context.Background(), tt.id, tt.tenantID)
			assert.Equal(t, tt.wantErr, err != nil)
			repo.AssertExpectations(t)
		})
	}
}
