package service

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/internal/model"
)

// MockUserRepositoryForTenant mock untuk UserRepository interface (digunakan di tenant service test)
type MockUserRepositoryForTenant struct {
	mock.Mock
}

func (m *MockUserRepositoryForTenant) FindByPhone(phone string) (*model.User, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepositoryForTenant) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepositoryForTenant) Create(user *model.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepositoryForTenant) UpdateLastLogin(userID string) error {
	return m.Called(userID).Error(0)
}

func (m *MockUserRepositoryForTenant) UpdatePassword(userID, hashedPassword string) error {
	return m.Called(userID, hashedPassword).Error(0)
}

func (m *MockUserRepositoryForTenant) UpdateVerifiedAt(userID string) error {
	return m.Called(userID).Error(0)
}

func (m *MockUserRepositoryForTenant) Delete(userID string) error {
	return m.Called(userID).Error(0)
}

func (m *MockUserRepositoryForTenant) GetAll() ([]model.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

// MockTenantRepository mock untuk TenantRepository interface
type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) CreateTenantWithMember(ctx context.Context, tenant *model.Tenant, member *model.TenantMember) error {
	return m.Called(ctx, tenant, member).Error(0)
}

func (m *MockTenantRepository) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	return m.Called(ctx, tenant).Error(0)
}

func (m *MockTenantRepository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetTenantsByUserID(ctx context.Context, userID string) ([]*model.Tenant, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Tenant), args.Error(1)
}

func (m *MockTenantRepository) UpdateTenant(ctx context.Context, tenant *model.Tenant) error {
	return m.Called(ctx, tenant).Error(0)
}

func (m *MockTenantRepository) DeleteTenant(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockTenantRepository) CountTenantsByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTenantRepository) CreateTenantMember(ctx context.Context, member *model.TenantMember) error {
	return m.Called(ctx, member).Error(0)
}

func (m *MockTenantRepository) GetTenantMembers(ctx context.Context, tenantID string) ([]*model.TenantMember, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.TenantMember), args.Error(1)
}

func (m *MockTenantRepository) IsTenantMember(ctx context.Context, tenantID, userID string) (bool, error) {
	args := m.Called(ctx, tenantID, userID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTenantRepository) IsTenantAdmin(ctx context.Context, tenantID, userID string) (bool, error) {
	args := m.Called(ctx, tenantID, userID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTenantRepository) RemoveTenantMember(ctx context.Context, tenantID, userID string) error {
	return m.Called(ctx, tenantID, userID).Error(0)
}

func (m *MockTenantRepository) IsSlugTaken(ctx context.Context, slug, excludeID string) (bool, error) {
	args := m.Called(ctx, slug, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

// MockContactCategoryRepositoryForTenant mock untuk ContactCategoryRepository (digunakan di tenant service test)
type MockContactCategoryRepositoryForTenant struct {
	mock.Mock
}

func (m *MockContactCategoryRepositoryForTenant) Create(ctx context.Context, category *model.ContactCategory) error {
	return m.Called(ctx, category).Error(0)
}

func (m *MockContactCategoryRepositoryForTenant) FindByTenantID(ctx context.Context, tenantID string) ([]*model.ContactCategory, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ContactCategory), args.Error(1)
}

func (m *MockContactCategoryRepositoryForTenant) FindByID(ctx context.Context, id string) (*model.ContactCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContactCategory), args.Error(1)
}

func (m *MockContactCategoryRepositoryForTenant) Update(ctx context.Context, category *model.ContactCategory) error {
	return m.Called(ctx, category).Error(0)
}

func (m *MockContactCategoryRepositoryForTenant) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockNotificationServiceForTenant mock untuk NotificationService
type MockNotificationServiceForTenant struct {
	mock.Mock
}

func (m *MockNotificationServiceForTenant) Send(ctx context.Context, phone, message string) error {
	return m.Called(ctx, phone, message).Error(0)
}

// Test CreateTenant - first tenant (free)
func TestCreateTenant_FirstTenant(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	userID := "user-123"
	tenantName := "Kelompok Masjid Al-Ikhlas"

	mockRepo.On("CountTenantsByOwner", ctx, userID).Return(int64(0), nil)
	mockRepo.On("IsSlugTaken", ctx, mock.AnythingOfType("string"), "").Return(false, nil)
	mockRepo.On("CreateTenantWithMember", ctx,
		mock.MatchedBy(func(t *model.Tenant) bool {
			return t.Name == tenantName && t.Status == "free" && t.ExpiredAt == nil && t.OwnerID == userID
		}),
		mock.MatchedBy(func(m *model.TenantMember) bool {
			return m.UserID == userID && m.Role == "admin"
		}),
	).Return(nil)
	mockCategoryRepo.On("Create", ctx, mock.Anything).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	mockUserRepo.On("FindByID", userID).Return(&model.User{ID: userID, Phone: "+6281234567890"}, nil)
	mockNotif.On("Send", ctx, "+6281234567890", mock.MatchedBy(func(msg string) bool {
		return len(msg) > 0
	})).Return(nil)

	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.CreateTenant(ctx, userID, tenantName, nil)

	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, "free", tenant.Status)
	assert.Nil(t, tenant.ExpiredAt)
	assert.Equal(t, userID, tenant.OwnerID)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockNotif.AssertExpectations(t)
}

// Test CreateTenant - second tenant (paid)
func TestCreateTenant_SecondTenant(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	userID := "user-123"
	tenantName := "Kelompok Masjid Nurul Huda"

	mockRepo.On("CountTenantsByOwner", ctx, userID).Return(int64(1), nil)
	mockRepo.On("IsSlugTaken", ctx, mock.AnythingOfType("string"), "").Return(false, nil)
	mockRepo.On("CreateTenantWithMember", ctx,
		mock.MatchedBy(func(t *model.Tenant) bool {
			return t.Name == tenantName && t.Status == "paid" && t.ExpiredAt != nil && t.OwnerID == userID
		}),
		mock.MatchedBy(func(m *model.TenantMember) bool {
			return m.UserID == userID && m.Role == "admin"
		}),
	).Return(nil)
	mockCategoryRepo.On("Create", ctx, mock.Anything).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	mockUserRepo.On("FindByID", userID).Return(&model.User{ID: userID, Phone: "+6281234567890"}, nil)
	mockNotif.On("Send", ctx, "+6281234567890", mock.MatchedBy(func(msg string) bool {
		return len(msg) > 0
	})).Return(nil)

	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.CreateTenant(ctx, userID, tenantName, nil)

	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, "paid", tenant.Status)
	assert.NotNil(t, tenant.ExpiredAt)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockNotif.AssertExpectations(t)
}

// Test GetTenantByID
func TestGetTenantByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	expectedTenant := &model.Tenant{
		ID:      tenantID,
		Name:    "Kelompok Masjid Al-Ikhlas",
		Slug:    "kelompok-masjid-al-ikhlas",
		Status:  "free",
		OwnerID: "user-123",
	}

	mockRepo.On("GetTenantByID", ctx, tenantID).Return(expectedTenant, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.GetTenantByID(ctx, tenantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTenant, tenant)
	mockRepo.AssertExpectations(t)
}

// Test GetTenantByID - not found
func TestGetTenantByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-invalid"

	mockRepo.On("GetTenantByID", ctx, tenantID).Return(nil, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	_, err := service.GetTenantByID(ctx, tenantID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant tidak ditemukan")
	mockRepo.AssertExpectations(t)
}

// Test UpdateTenant - admin update
func TestUpdateTenant_AdminCanUpdate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	userID := "user-123"
	newName := "Kelompok Baru"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, userID).Return(true, nil)
	mockRepo.On("GetTenantByID", ctx, tenantID).Return(&model.Tenant{
		ID:      tenantID,
		Name:    "Kelompok Lama",
		OwnerID: userID,
	}, nil)
	mockRepo.On("IsSlugTaken", ctx, mock.AnythingOfType("string"), tenantID).Return(false, nil)
	mockRepo.On("UpdateTenant", ctx, mock.MatchedBy(func(t *model.Tenant) bool {
		return t.Name == newName
	})).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.UpdateTenant(ctx, tenantID, userID, newName, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, newName, tenant.Name)
	mockRepo.AssertExpectations(t)
}

// Test UpdateTenant - non-admin cannot update
func TestUpdateTenant_NonAdminCannotUpdate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	userID := "user-456"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, userID).Return(false, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	_, err := service.UpdateTenant(ctx, tenantID, userID, "New Name", nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user bukan admin")
	mockRepo.AssertExpectations(t)
}

// Test InviteTenantMember
func TestInviteTenantMember(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	invitedUserID := "user-456"
	invitedBy := "user-123"

	mockUserRepo := new(MockUserRepositoryForTenant)
	mockUserRepo.On("FindByID", invitedUserID).Return(&model.User{ID: invitedUserID, Phone: "+6289876543210"}, nil)
	mockRepo.On("IsTenantMember", ctx, tenantID, invitedUserID).Return(false, nil)
	mockRepo.On("CreateTenantMember", ctx, mock.MatchedBy(func(m *model.TenantMember) bool {
		return m.TenantID == tenantID && m.UserID == invitedUserID && m.Role == "guest" && m.InvitedBy != nil && *m.InvitedBy == invitedBy
	})).Run(func(args mock.Arguments) {
		member := args.Get(1).(*model.TenantMember)
		member.ID = "member-123"
	}).Return(nil)

	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	member, err := service.InviteTenantMember(ctx, tenantID, invitedUserID, invitedBy)

	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "guest", member.Role)
	assert.Equal(t, invitedUserID, member.UserID)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test InviteTenantMember - user already member
func TestInviteTenantMember_AlreadyMember(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	invitedUserID := "user-456"
	invitedBy := "user-123"

	mockUserRepo := new(MockUserRepositoryForTenant)
	mockUserRepo.On("FindByID", invitedUserID).Return(&model.User{ID: invitedUserID, Phone: "+6289876543210"}, nil)
	mockRepo.On("IsTenantMember", ctx, tenantID, invitedUserID).Return(true, nil)

	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	_, err := service.InviteTenantMember(ctx, tenantID, invitedUserID, invitedBy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sudah menjadi member")
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test InviteTenantMember - user tidak ditemukan
func TestInviteTenantMember_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	invitedUserID := "user-invalid"
	invitedBy := "user-123"

	mockUserRepo := new(MockUserRepositoryForTenant)
	mockUserRepo.On("FindByID", invitedUserID).Return(nil, nil)

	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	_, err := service.InviteTenantMember(ctx, tenantID, invitedUserID, invitedBy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user tidak ditemukan")
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test RemoveTenantMember - admin can remove
func TestRemoveTenantMember_AdminCanRemove(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	memberUserID := "user-456"
	requesterID := "user-123"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, requesterID).Return(true, nil)
	mockRepo.On("GetTenantByID", ctx, tenantID).Return(&model.Tenant{
		ID:      tenantID,
		OwnerID: requesterID, // owner adalah requester, bukan memberUserID
	}, nil)
	mockRepo.On("IsTenantMember", ctx, tenantID, memberUserID).Return(true, nil)
	mockRepo.On("RemoveTenantMember", ctx, tenantID, memberUserID).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	err := service.RemoveTenantMember(ctx, tenantID, memberUserID, requesterID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test RemoveTenantMember - non-admin cannot remove
func TestRemoveTenantMember_NonAdminCannotRemove(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	memberUserID := "user-456"
	requesterID := "user-789"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, requesterID).Return(false, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	err := service.RemoveTenantMember(ctx, tenantID, memberUserID, requesterID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user bukan admin")
	mockRepo.AssertExpectations(t)
}

// Test RemoveTenantMember - owner tidak bisa dihapus
func TestRemoveTenantMember_CannotRemoveOwner(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	ownerID := "user-123"
	requesterID := "user-123" // requester adalah admin sekaligus owner

	mockRepo.On("IsTenantAdmin", ctx, tenantID, requesterID).Return(true, nil)
	mockRepo.On("GetTenantByID", ctx, tenantID).Return(&model.Tenant{
		ID:      tenantID,
		OwnerID: ownerID,
	}, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	err := service.RemoveTenantMember(ctx, tenantID, ownerID, requesterID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "owner tenant tidak bisa dihapus")
	mockRepo.AssertExpectations(t)
}

// Test DeleteTenant - admin can delete
func TestDeleteTenant_AdminCanDelete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	userID := "user-123"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, userID).Return(true, nil)
	mockRepo.On("DeleteTenant", ctx, tenantID).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	err := service.DeleteTenant(ctx, tenantID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test DeleteTenant - non-admin cannot delete
func TestDeleteTenant_NonAdminCannotDelete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	userID := "user-456"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, userID).Return(false, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	err := service.DeleteTenant(ctx, tenantID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user bukan admin")
	mockRepo.AssertExpectations(t)
}

// Test GetUserTenants
func TestGetUserTenants(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	userID := "user-123"
	tenants := []*model.Tenant{
		{ID: "tenant-1", Name: "Tenant 1"},
		{ID: "tenant-2", Name: "Tenant 2"},
	}

	mockRepo.On("GetTenantsByUserID", ctx, userID).Return(tenants, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	result, err := service.GetUserTenants(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// Test GetTenantMembers
func TestGetTenantMembers(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	members := []*model.TenantMember{
		{ID: "member-1", UserID: "user-1", Role: "admin"},
		{ID: "member-2", UserID: "user-2", Role: "guest"},
	}

	mockRepo.On("GetTenantMembers", ctx, tenantID).Return(members, nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	result, err := service.GetTenantMembers(ctx, tenantID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// Test CreateTenant with transaction rollback when CreateTenantWithMember fails
func TestCreateTenant_TransactionRollback(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	userID := "user-123"
	tenantName := "Kelompok Masjid"

	mockRepo.On("CountTenantsByOwner", ctx, userID).Return(int64(0), nil)
	mockRepo.On("IsSlugTaken", ctx, mock.AnythingOfType("string"), "").Return(false, nil)
	// Transaksi gagal — tenant tidak tersimpan (rollback)
	mockRepo.On("CreateTenantWithMember", ctx, mock.Anything, mock.Anything).Return(errors.New("transaction failed"))

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.CreateTenant(ctx, userID, tenantName, nil)

	assert.Error(t, err)
	assert.Nil(t, tenant)
	mockRepo.AssertExpectations(t)
}

// Test CreateTenant - slug collision menghasilkan suffix
func TestCreateTenant_SlugCollision(t *testing.T) {	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	userID := "user-123"
	tenantName := "Masjid Al Ikhlas"
	expectedBaseSlug := "masjid-al-ikhlas"

	mockRepo.On("CountTenantsByOwner", ctx, userID).Return(int64(0), nil)
	// Slug pertama sudah dipakai, slug dengan suffix -1 tersedia
	mockRepo.On("IsSlugTaken", ctx, expectedBaseSlug, "").Return(true, nil)
	mockRepo.On("IsSlugTaken", ctx, expectedBaseSlug+"-1", "").Return(false, nil)
	mockRepo.On("CreateTenantWithMember", ctx,
		mock.MatchedBy(func(t *model.Tenant) bool {
			return t.Slug == expectedBaseSlug+"-1"
		}),
		mock.Anything,
	).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	mockUserRepo.On("FindByID", userID).Return(&model.User{ID: userID, Phone: "+6281234567890"}, nil)
	mockNotif.On("Send", ctx, "+6281234567890", mock.Anything).Return(nil)
	mockCategoryRepo.On("Create", ctx, mock.Anything).Return(nil)

	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.CreateTenant(ctx, userID, tenantName, nil)

	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, expectedBaseSlug+"-1", tenant.Slug)
	mockRepo.AssertExpectations(t)
}

// Test UpdateTenant - dengan address dan description diisi
func TestUpdateTenant_WithAddressAndDescription(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	userID := "user-123"
	newName := "Kelompok Baru"
	addr := "Jl. Sudirman No. 1, Jakarta"
	desc := "Deskripsi kelompok qurban"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, userID).Return(true, nil)
	mockRepo.On("GetTenantByID", ctx, tenantID).Return(&model.Tenant{
		ID:      tenantID,
		Name:    "Kelompok Lama",
		OwnerID: userID,
	}, nil)
	mockRepo.On("IsSlugTaken", ctx, mock.AnythingOfType("string"), tenantID).Return(false, nil)
	mockRepo.On("UpdateTenant", ctx, mock.MatchedBy(func(t *model.Tenant) bool {
		return t.Name == newName && t.Address != nil && *t.Address == addr && t.Description != nil && *t.Description == desc
	})).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.UpdateTenant(ctx, tenantID, userID, newName, &addr, &desc, nil)

	assert.NoError(t, err)
	assert.Equal(t, newName, tenant.Name)
	assert.NotNil(t, tenant.Address)
	assert.Equal(t, addr, *tenant.Address)
	assert.NotNil(t, tenant.Description)
	assert.Equal(t, desc, *tenant.Description)
	mockRepo.AssertExpectations(t)
}

// Test UpdateTenant - address dan description nil (opsional tidak diisi)
func TestUpdateTenant_WithNilOptionalFields(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepository)
	mockNotif := new(MockNotificationServiceForTenant)
	mockCategoryRepo := new(MockContactCategoryRepositoryForTenant)
	logger := log.New(io.Discard, "", 0)

	tenantID := "tenant-123"
	userID := "user-123"
	newName := "Kelompok Diperbarui"

	mockRepo.On("IsTenantAdmin", ctx, tenantID, userID).Return(true, nil)
	mockRepo.On("GetTenantByID", ctx, tenantID).Return(&model.Tenant{
		ID:      tenantID,
		Name:    "Kelompok Lama",
		OwnerID: userID,
	}, nil)
	mockRepo.On("IsSlugTaken", ctx, mock.AnythingOfType("string"), tenantID).Return(false, nil)
	mockRepo.On("UpdateTenant", ctx, mock.MatchedBy(func(t *model.Tenant) bool {
		return t.Name == newName && t.Address == nil && t.Description == nil
	})).Return(nil)

	mockUserRepo := new(MockUserRepositoryForTenant)
	service := NewTenantService(mockRepo, mockUserRepo, mockNotif, mockCategoryRepo, logger)
	tenant, err := service.UpdateTenant(ctx, tenantID, userID, newName, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, newName, tenant.Name)
	assert.Nil(t, tenant.Address)
	assert.Nil(t, tenant.Description)
	mockRepo.AssertExpectations(t)
}
