package repository

import (
	"context"
	"fmt"

	"github.com/username/qurban-app/internal/model"
	"gorm.io/gorm"
)

// ContactCategoryRepository interface untuk operasi contact category
type ContactCategoryRepository interface {
	Create(ctx context.Context, category *model.ContactCategory) error
	FindByTenantID(ctx context.Context, tenantID string) ([]*model.ContactCategory, error)
	FindByID(ctx context.Context, id string) (*model.ContactCategory, error)
	Update(ctx context.Context, category *model.ContactCategory) error
	Delete(ctx context.Context, id string) error
}

// contactCategoryRepository implementasi ContactCategoryRepository
type contactCategoryRepository struct {
	db *gorm.DB
}

// NewContactCategoryRepository membuat instance baru ContactCategoryRepository
func NewContactCategoryRepository(db *gorm.DB) ContactCategoryRepository {
	return &contactCategoryRepository{db: db}
}

// Create menyimpan contact category baru ke database
func (r *contactCategoryRepository) Create(ctx context.Context, category *model.ContactCategory) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("Create: %w", err)
	}
	return nil
}

// FindByTenantID mengambil semua contact category berdasarkan tenant ID
func (r *contactCategoryRepository) FindByTenantID(ctx context.Context, tenantID string) ([]*model.ContactCategory, error) {
	var categories []*model.ContactCategory
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("FindByTenantID: %w", err)
	}
	return categories, nil
}

// FindByID mencari contact category berdasarkan ID
func (r *contactCategoryRepository) FindByID(ctx context.Context, id string) (*model.ContactCategory, error) {
	var category model.ContactCategory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return &category, nil
}

// Update mengupdate data contact category
func (r *contactCategoryRepository) Update(ctx context.Context, category *model.ContactCategory) error {
	if err := r.db.WithContext(ctx).Model(category).Updates(category).Error; err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	return nil
}

// Delete menghapus contact category secara soft delete
func (r *contactCategoryRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.ContactCategory{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}

// ContactRepository interface untuk operasi contact
type ContactRepository interface {
	Create(ctx context.Context, contact *model.Contact) error
	FindByTenantID(ctx context.Context, tenantID string) ([]*model.Contact, error)
	FindByID(ctx context.Context, id string) (*model.Contact, error)
	Update(ctx context.Context, contact *model.Contact) error
	Delete(ctx context.Context, id string) error
}

// contactRepository implementasi ContactRepository
type contactRepository struct {
	db *gorm.DB
}

// NewContactRepository membuat instance baru ContactRepository
func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

// Create menyimpan contact baru ke database
func (r *contactRepository) Create(ctx context.Context, contact *model.Contact) error {
	if err := r.db.WithContext(ctx).Create(contact).Error; err != nil {
		return fmt.Errorf("Create: %w", err)
	}
	return nil
}

// FindByTenantID mengambil semua contact berdasarkan tenant ID
func (r *contactRepository) FindByTenantID(ctx context.Context, tenantID string) ([]*model.Contact, error) {
	var contacts []*model.Contact
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&contacts).Error; err != nil {
		return nil, fmt.Errorf("FindByTenantID: %w", err)
	}
	return contacts, nil
}

// FindByID mencari contact berdasarkan ID
func (r *contactRepository) FindByID(ctx context.Context, id string) (*model.Contact, error) {
	var contact model.Contact
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&contact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return &contact, nil
}

// Update mengupdate data contact (menggunakan Save agar field bool zero value juga tersimpan)
func (r *contactRepository) Update(ctx context.Context, contact *model.Contact) error {
	if err := r.db.WithContext(ctx).Save(contact).Error; err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	return nil
}

// Delete menghapus contact secara soft delete
func (r *contactRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.Contact{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}
