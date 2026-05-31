package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/internal/repository"
)

// TenantService interface untuk operasi tenant
type TenantService interface {
	CreateTenant(ctx context.Context, userID, name string, logo *string) (*model.Tenant, error)
	GetTenantByID(ctx context.Context, id string) (*model.Tenant, error)
	GetUserTenants(ctx context.Context, userID string) ([]*model.Tenant, error)
	UpdateTenant(ctx context.Context, tenantID, userID, name string, address, description, logo *string) (*model.Tenant, error)
	DeleteTenant(ctx context.Context, tenantID, userID string) error
	InviteTenantMember(ctx context.Context, tenantID, invitedUserID, invitedBy string) (*model.TenantMember, error)
	GetTenantMembers(ctx context.Context, tenantID string) ([]*model.TenantMember, error)
	RemoveTenantMember(ctx context.Context, tenantID, userID, requesterID string) error
}

// tenantService implementasi TenantService
type tenantService struct {
	tenantRepo          repository.TenantRepository
	userRepo            repository.UserRepository
	notificationService NotificationService
	contactCategoryRepo repository.ContactCategoryRepository
	logger              *log.Logger
}

// NewTenantService membuat instance baru TenantService
func NewTenantService(tenantRepo repository.TenantRepository, userRepo repository.UserRepository, notifService NotificationService, contactCategoryRepo repository.ContactCategoryRepository, logger *log.Logger) TenantService {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &tenantService{
		tenantRepo:          tenantRepo,
		userRepo:            userRepo,
		notificationService: notifService,
		contactCategoryRepo: contactCategoryRepo,
		logger:              logger,
	}
}

// CreateTenant membuat tenant baru dengan logika business
func (s *tenantService) CreateTenant(ctx context.Context, userID, name string, logo *string) (*model.Tenant, error) {
	if userID == "" {
		return nil, fmt.Errorf("CreateTenant: userID tidak boleh kosong")
	}
	if name == "" {
		return nil, fmt.Errorf("CreateTenant: name tidak boleh kosong")
	}

	// Cek berapa banyak tenant yang sudah dimiliki user
	count, err := s.tenantRepo.CountTenantsByOwner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("CreateTenant: %w", err)
	}

	// Tentukan status dan expired_at berdasarkan count
	status := model.TenantStatusFree
	var expiredAt *time.Time
	if count > 0 {
		status = model.TenantStatusPaid
		exp := time.Now().AddDate(0, 0, 14) // +14 hari
		expiredAt = &exp
	}

	// Generate slug unik dari name
	slug, err := s.generateUniqueSlug(ctx, name, "")
	if err != nil {
		return nil, fmt.Errorf("CreateTenant: %w", err)
	}

	tenant := &model.Tenant{
		ID:        uuid.New().String(),
		Name:      name,
		Slug:      slug,
		Status:    string(status),
		Logo:      logo,
		ExpiredAt: expiredAt,
		OwnerID:   userID,
	}

	// Tambahkan user sebagai admin member otomatis
	member := &model.TenantMember{
		ID:       uuid.New().String(),
		TenantID: tenant.ID,
		UserID:   userID,
		Role:     string(model.TenantMemberRoleAdmin),
		JoinedAt: time.Now(),
	}

	// Buat tenant dan member dalam satu transaksi — rollback jika salah satu gagal
	if err := s.tenantRepo.CreateTenantWithMember(ctx, tenant, member); err != nil {
		return nil, fmt.Errorf("CreateTenant: %w", err)
	}

	// Buat kategori kontak default untuk tenant baru
	defaultCategory := &model.ContactCategory{
		ID:       uuid.New().String(),
		TenantID: tenant.ID,
		Name:     "umum",
	}
	if err := s.contactCategoryRepo.Create(ctx, defaultCategory); err != nil {
		s.logger.Printf("CreateTenant: gagal buat kategori default: %v\n", err)
		// tidak return error, lanjut saja
	}

	// Kirim notifikasi WhatsApp ke pembuat tenant
	s.sendTenantCreatedNotification(ctx, userID, tenant)

	return tenant, nil
}

// sendTenantCreatedNotification mengirim notifikasi WhatsApp setelah tenant berhasil dibuat
func (s *tenantService) sendTenantCreatedNotification(ctx context.Context, userID string, tenant *model.Tenant) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		s.logger.Printf("sendTenantCreatedNotification: gagal ambil data user %s: %v\n", userID, err)
		return
	}

	statusText := "Gratis"
	if tenant.Status == string(model.TenantStatusPaid) {
		statusText = "Berbayar"
	}

	expiredText := "Tidak ada"
	if tenant.ExpiredAt != nil {
		expiredText = tenant.ExpiredAt.Format("02 January 2006")
	}

	message := fmt.Sprintf(
		"Organisasi Anda berhasil dibuat!\nNama: %s\nStatus: %s\nMasa berlaku: %s",
		tenant.Name,
		statusText,
		expiredText,
	)

	if err := s.notificationService.Send(ctx, user.Phone, message); err != nil {
		s.logger.Printf("sendTenantCreatedNotification: gagal kirim notifikasi WhatsApp ke %s: %v\n", user.Phone, err)
	}
}

// GetTenantByID mengambil detail tenant berdasarkan ID
func (s *tenantService) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	if id == "" {
		return nil, fmt.Errorf("GetTenantByID: id tidak boleh kosong")
	}

	tenant, err := s.tenantRepo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetTenantByID: %w", err)
	}

	if tenant == nil {
		return nil, fmt.Errorf("GetTenantByID: tenant tidak ditemukan")
	}

	return tenant, nil
}

// GetUserTenants mengambil semua tenant milik user
func (s *tenantService) GetUserTenants(ctx context.Context, userID string) ([]*model.Tenant, error) {
	if userID == "" {
		return nil, fmt.Errorf("GetUserTenants: userID tidak boleh kosong")
	}

	tenants, err := s.tenantRepo.GetTenantsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserTenants: %w", err)
	}

	if tenants == nil {
		tenants = []*model.Tenant{}
	}

	return tenants, nil
}

// UpdateTenant mengupdate data tenant (hanya admin yang bisa)
func (s *tenantService) UpdateTenant(ctx context.Context, tenantID, userID, name string, address, description, logo *string) (*model.Tenant, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("UpdateTenant: tenantID tidak boleh kosong")
	}
	if userID == "" {
		return nil, fmt.Errorf("UpdateTenant: userID tidak boleh kosong")
	}
	if name == "" {
		return nil, fmt.Errorf("UpdateTenant: name tidak boleh kosong")
	}

	// Cek apakah user adalah admin di tenant
	isAdmin, err := s.tenantRepo.IsTenantAdmin(ctx, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTenant: %w", err)
	}
	if !isAdmin {
		return nil, fmt.Errorf("UpdateTenant: user bukan admin di tenant ini")
	}

	// Ambil tenant lama
	tenant, err := s.tenantRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("UpdateTenant: tenant tidak ditemukan")
	}

	// Update nama dan slug unik (excludeID = tenantID agar slug milik tenant ini tidak dianggap konflik)
	newSlug, err := s.generateUniqueSlug(ctx, name, tenantID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTenant: %w", err)
	}
	tenant.Name = name
	tenant.Slug = newSlug
	tenant.Address = address
	tenant.Description = description
	tenant.Logo = logo

	if err := s.tenantRepo.UpdateTenant(ctx, tenant); err != nil {
		return nil, fmt.Errorf("UpdateTenant: %w", err)
	}

	return tenant, nil
}

// DeleteTenant menghapus tenant (hanya admin yang bisa)
func (s *tenantService) DeleteTenant(ctx context.Context, tenantID, userID string) error {
	if tenantID == "" {
		return fmt.Errorf("DeleteTenant: tenantID tidak boleh kosong")
	}
	if userID == "" {
		return fmt.Errorf("DeleteTenant: userID tidak boleh kosong")
	}

	// Cek apakah user adalah admin di tenant
	isAdmin, err := s.tenantRepo.IsTenantAdmin(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("DeleteTenant: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("DeleteTenant: user bukan admin di tenant ini")
	}

	if err := s.tenantRepo.DeleteTenant(ctx, tenantID); err != nil {
		return fmt.Errorf("DeleteTenant: %w", err)
	}

	return nil
}

// InviteTenantMember mengundang user baru sebagai member di tenant
func (s *tenantService) InviteTenantMember(ctx context.Context, tenantID, invitedUserID, invitedBy string) (*model.TenantMember, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("InviteTenantMember: tenantID tidak boleh kosong")
	}
	if invitedUserID == "" {
		return nil, fmt.Errorf("InviteTenantMember: invitedUserID tidak boleh kosong")
	}
	if invitedBy == "" {
		return nil, fmt.Errorf("InviteTenantMember: invitedBy tidak boleh kosong")
	}

	// Validasi bahwa user yang diundang benar-benar ada di sistem
	user, err := s.userRepo.FindByID(invitedUserID)
	if err != nil {
		return nil, fmt.Errorf("InviteTenantMember: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("InviteTenantMember: user tidak ditemukan")
	}

	// Cek apakah user yang diundang sudah menjadi member
	isMember, err := s.tenantRepo.IsTenantMember(ctx, tenantID, invitedUserID)
	if err != nil {
		return nil, fmt.Errorf("InviteTenantMember: %w", err)
	}
	if isMember {
		return nil, fmt.Errorf("InviteTenantMember: user sudah menjadi member di tenant ini")
	}

	// Buat member baru dengan role guest
	member := &model.TenantMember{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		UserID:    invitedUserID,
		Role:      string(model.TenantMemberRoleGuest),
		InvitedBy: &invitedBy,
		JoinedAt:  time.Now(),
	}

	if err := s.tenantRepo.CreateTenantMember(ctx, member); err != nil {
		return nil, fmt.Errorf("InviteTenantMember: %w", err)
	}

	return member, nil
}

// GetTenantMembers mengambil semua member di tenant
func (s *tenantService) GetTenantMembers(ctx context.Context, tenantID string) ([]*model.TenantMember, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("GetTenantMembers: tenantID tidak boleh kosong")
	}

	members, err := s.tenantRepo.GetTenantMembers(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("GetTenantMembers: %w", err)
	}

	if members == nil {
		members = []*model.TenantMember{}
	}

	return members, nil
}

// RemoveTenantMember menghapus member dari tenant (hanya admin yang bisa)
func (s *tenantService) RemoveTenantMember(ctx context.Context, tenantID, userID, requesterID string) error {
	if tenantID == "" {
		return fmt.Errorf("RemoveTenantMember: tenantID tidak boleh kosong")
	}
	if userID == "" {
		return fmt.Errorf("RemoveTenantMember: userID tidak boleh kosong")
	}
	if requesterID == "" {
		return fmt.Errorf("RemoveTenantMember: requesterID tidak boleh kosong")
	}

	// Cek apakah requester adalah admin di tenant
	isAdmin, err := s.tenantRepo.IsTenantAdmin(ctx, tenantID, requesterID)
	if err != nil {
		return fmt.Errorf("RemoveTenantMember: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("RemoveTenantMember: user bukan admin di tenant ini")
	}

	// Cek apakah user yang akan dihapus adalah owner tenant
	tenant, err := s.tenantRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("RemoveTenantMember: %w", err)
	}
	if tenant == nil {
		return fmt.Errorf("RemoveTenantMember: tenant tidak ditemukan")
	}
	if tenant.OwnerID == userID {
		return fmt.Errorf("RemoveTenantMember: owner tenant tidak bisa dihapus dari member")
	}

	// Cek apakah user yang dihapus adalah member
	isMember, err := s.tenantRepo.IsTenantMember(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("RemoveTenantMember: %w", err)
	}
	if !isMember {
		return fmt.Errorf("RemoveTenantMember: user bukan member di tenant ini")
	}

	if err := s.tenantRepo.RemoveTenantMember(ctx, tenantID, userID); err != nil {
		return fmt.Errorf("RemoveTenantMember: %w", err)
	}

	return nil
}

// generateUniqueSlug membuat slug unik dari nama. excludeID diisi jika ingin mengecualikan tenant tertentu (untuk update).
func (s *tenantService) generateUniqueSlug(ctx context.Context, name, excludeID string) (string, error) {
	baseSlug := generateSlug(name)
	slug := baseSlug
	for i := 1; i <= 100; i++ {
		taken, err := s.tenantRepo.IsSlugTaken(ctx, slug, excludeID)
		if err != nil {
			return "", fmt.Errorf("generateUniqueSlug: %w", err)
		}
		if !taken {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}
	return "", fmt.Errorf("generateUniqueSlug: tidak bisa membuat slug unik untuk nama %q", name)
}

// generateSlug membuat slug dari nama
func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}
