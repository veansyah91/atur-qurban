package repository

import (
	"context"
	"fmt"

	"github.com/username/qurban-app/internal/model"
	"gorm.io/gorm"
)

// TenantRepository interface untuk operasi tenant
type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *model.Tenant) error
	CreateTenantWithMember(ctx context.Context, tenant *model.Tenant, member *model.TenantMember) error
	GetTenantByID(ctx context.Context, id string) (*model.Tenant, error)
	GetTenantsByUserID(ctx context.Context, userID string) ([]*model.Tenant, error)
	UpdateTenant(ctx context.Context, tenant *model.Tenant) error
	DeleteTenant(ctx context.Context, id string) error
	CountTenantsByOwner(ctx context.Context, ownerID string) (int64, error)
	IsSlugTaken(ctx context.Context, slug, excludeID string) (bool, error)

	CreateTenantMember(ctx context.Context, member *model.TenantMember) error
	GetTenantMembers(ctx context.Context, tenantID string) ([]*model.TenantMember, error)
	IsTenantMember(ctx context.Context, tenantID, userID string) (bool, error)
	IsTenantAdmin(ctx context.Context, tenantID, userID string) (bool, error)
	RemoveTenantMember(ctx context.Context, tenantID, userID string) error
}

// tenantRepository implementasi TenantRepository
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository membuat instance baru TenantRepository
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

// CreateTenant menyimpan tenant baru ke database
func (r *tenantRepository) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	if err := r.db.WithContext(ctx).Create(tenant).Error; err != nil {
		return fmt.Errorf("CreateTenant: %w", err)
	}
	return nil
}

// GetTenantByID mencari tenant berdasarkan ID
func (r *tenantRepository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("GetTenantByID: %w", err)
	}
	return &tenant, nil
}

// GetTenantsByUserID mengambil semua tenant dimana user adalah member
func (r *tenantRepository) GetTenantsByUserID(ctx context.Context, userID string) ([]*model.Tenant, error) {
	var tenants []*model.Tenant
	if err := r.db.WithContext(ctx).
		Joins("INNER JOIN tenant_members ON tenant_members.tenant_id = tenants.id").
		Where("tenant_members.user_id = ? AND tenant_members.deleted_at IS NULL", userID).
		Find(&tenants).Error; err != nil {
		return nil, fmt.Errorf("GetTenantsByUserID: %w", err)
	}
	return tenants, nil
}

// UpdateTenant mengupdate data tenant
func (r *tenantRepository) UpdateTenant(ctx context.Context, tenant *model.Tenant) error {
	if err := r.db.WithContext(ctx).Model(tenant).Updates(tenant).Error; err != nil {
		return fmt.Errorf("UpdateTenant: %w", err)
	}
	return nil
}

// DeleteTenant menghapus tenant secara soft delete
func (r *tenantRepository) DeleteTenant(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.Tenant{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("DeleteTenant: %w", err)
	}
	return nil
}

// CountTenantsByOwner menghitung jumlah tenant milik owner
func (r *tenantRepository) CountTenantsByOwner(ctx context.Context, ownerID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Tenant{}).Where("owner_id = ?", ownerID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("CountTenantsByOwner: %w", err)
	}
	return count, nil
}

// IsSlugTaken mengecek apakah slug sudah dipakai oleh tenant lain. excludeID dikosongkan jika tidak ada pengecualian.
func (r *tenantRepository) IsSlugTaken(ctx context.Context, slug, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Tenant{}).Where("slug = ?", slug)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("IsSlugTaken: %w", err)
	}
	return count > 0, nil
}

// CreateTenantWithMember membuat tenant dan member awal dalam satu transaksi database
func (r *tenantRepository) CreateTenantWithMember(ctx context.Context, tenant *model.Tenant, member *model.TenantMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			return fmt.Errorf("CreateTenantWithMember: buat tenant: %w", err)
		}
		member.TenantID = tenant.ID
		if err := tx.Create(member).Error; err != nil {
			return fmt.Errorf("CreateTenantWithMember: buat member: %w", err)
		}
		return nil
	})
}

// CreateTenantMember menambahkan member baru ke tenant
func (r *tenantRepository) CreateTenantMember(ctx context.Context, member *model.TenantMember) error {
	if err := r.db.WithContext(ctx).Create(member).Error; err != nil {
		return fmt.Errorf("CreateTenantMember: %w", err)
	}
	return nil
}

// GetTenantMembers mengambil semua member dari tenant
func (r *tenantRepository) GetTenantMembers(ctx context.Context, tenantID string) ([]*model.TenantMember, error) {
	var members []*model.TenantMember
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&members).Error; err != nil {
		return nil, fmt.Errorf("GetTenantMembers: %w", err)
	}
	return members, nil
}

// IsTenantMember mengecek apakah user adalah member di tenant
func (r *tenantRepository) IsTenantMember(ctx context.Context, tenantID, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.TenantMember{}).
		Where("tenant_id = ? AND user_id = ? AND deleted_at IS NULL", tenantID, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("IsTenantMember: %w", err)
	}
	return count > 0, nil
}

// IsTenantAdmin mengecek apakah user adalah admin di tenant
func (r *tenantRepository) IsTenantAdmin(ctx context.Context, tenantID, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.TenantMember{}).
		Where("tenant_id = ? AND user_id = ? AND role = ? AND deleted_at IS NULL", tenantID, userID, model.TenantMemberRoleAdmin).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("IsTenantAdmin: %w", err)
	}
	return count > 0, nil
}

// RemoveTenantMember menghapus member dari tenant secara soft delete
func (r *tenantRepository) RemoveTenantMember(ctx context.Context, tenantID, userID string) error {
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Delete(&model.TenantMember{}).Error; err != nil {
		return fmt.Errorf("RemoveTenantMember: %w", err)
	}
	return nil
}
