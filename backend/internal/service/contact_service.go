package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/internal/repository"
)

// ContactCategoryService interface untuk operasi contact category
type ContactCategoryService interface {
	Create(ctx context.Context, tenantID, name string) (*model.ContactCategory, error)
	GetByTenant(ctx context.Context, tenantID string) ([]*model.ContactCategory, error)
	GetByID(ctx context.Context, id, tenantID string) (*model.ContactCategory, error)
	Update(ctx context.Context, id, tenantID, name string) (*model.ContactCategory, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// contactCategoryService implementasi ContactCategoryService
type contactCategoryService struct {
	repo repository.ContactCategoryRepository
}

// NewContactCategoryService membuat instance baru ContactCategoryService
func NewContactCategoryService(repo repository.ContactCategoryRepository) ContactCategoryService {
	return &contactCategoryService{repo: repo}
}

// Create membuat contact category baru
func (s *contactCategoryService) Create(ctx context.Context, tenantID, name string) (*model.ContactCategory, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("Create: tenantID tidak boleh kosong")
	}
	if name == "" {
		return nil, fmt.Errorf("Create: name tidak boleh kosong")
	}

	category := &model.ContactCategory{
		ID:       uuid.New().String(),
		TenantID: tenantID,
		Name:     name,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	return category, nil
}

// GetByTenant mengambil semua contact category milik tenant
func (s *contactCategoryService) GetByTenant(ctx context.Context, tenantID string) ([]*model.ContactCategory, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("GetByTenant: tenantID tidak boleh kosong")
	}

	categories, err := s.repo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("GetByTenant: %w", err)
	}

	if categories == nil {
		categories = []*model.ContactCategory{}
	}

	return categories, nil
}

// GetByID mengambil detail contact category berdasarkan ID dan memverifikasi kepemilikan tenant
func (s *contactCategoryService) GetByID(ctx context.Context, id, tenantID string) (*model.ContactCategory, error) {
	if id == "" {
		return nil, fmt.Errorf("GetByID: id tidak boleh kosong")
	}

	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}

	if category == nil || category.TenantID != tenantID {
		return nil, fmt.Errorf("GetByID: contact category tidak ditemukan")
	}

	return category, nil
}

// Update mengupdate nama contact category
func (s *contactCategoryService) Update(ctx context.Context, id, tenantID, name string) (*model.ContactCategory, error) {
	if id == "" {
		return nil, fmt.Errorf("Update: id tidak boleh kosong")
	}

	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	if category == nil || category.TenantID != tenantID {
		return nil, fmt.Errorf("Update: contact category tidak ditemukan")
	}

	category.Name = name

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	return category, nil
}

// Delete menghapus contact category secara soft delete
func (s *contactCategoryService) Delete(ctx context.Context, id, tenantID string) error {
	if id == "" {
		return fmt.Errorf("Delete: id tidak boleh kosong")
	}

	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	if category == nil || category.TenantID != tenantID {
		return fmt.Errorf("Delete: contact category tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	return nil
}

// ContactService interface untuk operasi contact
type ContactService interface {
	Create(ctx context.Context, tenantID, categoryID, name string, address, phone *string) (*model.Contact, error)
	GetByTenant(ctx context.Context, tenantID string) ([]*model.Contact, error)
	GetByID(ctx context.Context, id, tenantID string) (*model.Contact, error)
	Update(ctx context.Context, id, tenantID, categoryID, name string, isActive bool, address, phone *string) (*model.Contact, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// contactService implementasi ContactService
type contactService struct {
	repo         repository.ContactRepository
	categoryRepo repository.ContactCategoryRepository
}

// NewContactService membuat instance baru ContactService
func NewContactService(repo repository.ContactRepository, categoryRepo repository.ContactCategoryRepository) ContactService {
	return &contactService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// Create membuat contact baru setelah memverifikasi kategori
func (s *contactService) Create(ctx context.Context, tenantID, categoryID, name string, address, phone *string) (*model.Contact, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("Create: tenantID tidak boleh kosong")
	}
	if categoryID == "" {
		return nil, fmt.Errorf("Create: categoryID tidak boleh kosong")
	}
	if name == "" {
		return nil, fmt.Errorf("Create: name tidak boleh kosong")
	}

	// Verifikasi bahwa kategori ada dan milik tenant ini
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	if category == nil || category.TenantID != tenantID {
		return nil, fmt.Errorf("Create: kategori tidak ditemukan atau bukan milik tenant ini")
	}

	contact := &model.Contact{
		ID:                uuid.New().String(),
		TenantID:          tenantID,
		ContactCategoryID: categoryID,
		Name:              name,
		IsActive:          true,
		Address:           address,
		Phone:             phone,
	}

	if err := s.repo.Create(ctx, contact); err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	return contact, nil
}

// GetByTenant mengambil semua contact milik tenant
func (s *contactService) GetByTenant(ctx context.Context, tenantID string) ([]*model.Contact, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("GetByTenant: tenantID tidak boleh kosong")
	}

	contacts, err := s.repo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("GetByTenant: %w", err)
	}

	if contacts == nil {
		contacts = []*model.Contact{}
	}

	return contacts, nil
}

// GetByID mengambil detail contact berdasarkan ID dan memverifikasi kepemilikan tenant
func (s *contactService) GetByID(ctx context.Context, id, tenantID string) (*model.Contact, error) {
	contact, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}

	if contact == nil || contact.TenantID != tenantID {
		return nil, fmt.Errorf("GetByID: contact tidak ditemukan")
	}

	return contact, nil
}

// Update mengupdate data contact
func (s *contactService) Update(ctx context.Context, id, tenantID, categoryID, name string, isActive bool, address, phone *string) (*model.Contact, error) {
	contact, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	if contact == nil || contact.TenantID != tenantID {
		return nil, fmt.Errorf("Update: contact tidak ditemukan")
	}

	// Verifikasi kategori baru ada dan milik tenant ini
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}
	if category == nil || category.TenantID != tenantID {
		return nil, fmt.Errorf("Update: kategori tidak ditemukan atau bukan milik tenant ini")
	}

	contact.ContactCategoryID = categoryID
	contact.Name = name
	contact.IsActive = isActive
	contact.Address = address
	contact.Phone = phone

	if err := s.repo.Update(ctx, contact); err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	return contact, nil
}

// Delete menghapus contact secara soft delete
func (s *contactService) Delete(ctx context.Context, id, tenantID string) error {
	contact, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	if contact == nil || contact.TenantID != tenantID {
		return fmt.Errorf("Delete: contact tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	return nil
}
