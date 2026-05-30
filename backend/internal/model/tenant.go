package model

import (
	"time"

	"gorm.io/gorm"
)

// TenantStatus enum untuk status tenant
type TenantStatus string

const (
	TenantStatusFree TenantStatus = "free"
	TenantStatusPaid TenantStatus = "paid"
)

// TenantMemberRole enum untuk role di dalam tenant
type TenantMemberRole string

const (
	TenantMemberRoleAdmin TenantMemberRole = "admin"
	TenantMemberRoleGuest TenantMemberRole = "guest"
)

// Tenant model untuk tabel tenants
type Tenant struct {
	ID          string         `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Status      string         `gorm:"type:varchar(50);not null;default:'free'" json:"status"`
	Address     *string        `gorm:"type:varchar(500);default:NULL" json:"address"`
	Description *string        `gorm:"type:text;default:NULL" json:"description"`
	ExpiredAt   *time.Time     `gorm:"type:timestamp;default:NULL" json:"expired_at"`
	OwnerID     string         `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (Tenant) TableName() string {
	return "tenants"
}

// TenantMember model untuk tabel tenant_members
type TenantMember struct {
	ID        string     `gorm:"primaryKey;type:uuid" json:"id"`
	TenantID  string     `gorm:"type:uuid;not null" json:"tenant_id"`
	UserID    string     `gorm:"type:uuid;not null" json:"user_id"`
	Role      string     `gorm:"type:varchar(50);not null;default:'guest'" json:"role"`
	InvitedBy *string    `gorm:"type:uuid;default:NULL" json:"invited_by"`
	JoinedAt  time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"joined_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (TenantMember) TableName() string {
	return "tenant_members"
}

// TenantResponse untuk response API tenant
type TenantResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	Address     *string    `json:"address"`
	Description *string    `json:"description"`
	ExpiredAt   *time.Time `json:"expired_at"`
	OwnerID     string     `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse convert Tenant ke TenantResponse
func (t *Tenant) ToResponse() *TenantResponse {
	return &TenantResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Status:      t.Status,
		Address:     t.Address,
		Description: t.Description,
		ExpiredAt:   t.ExpiredAt,
		OwnerID:     t.OwnerID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// TenantMemberResponse untuk response API tenant member
type TenantMemberResponse struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	UserID    string     `json:"user_id"`
	Role      string     `json:"role"`
	InvitedBy *string    `json:"invited_by"`
	JoinedAt  time.Time  `json:"joined_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse convert TenantMember ke TenantMemberResponse
func (tm *TenantMember) ToResponse() *TenantMemberResponse {
	return &TenantMemberResponse{
		ID:        tm.ID,
		TenantID:  tm.TenantID,
		UserID:    tm.UserID,
		Role:      tm.Role,
		InvitedBy: tm.InvitedBy,
		JoinedAt:  tm.JoinedAt,
		CreatedAt: tm.CreatedAt,
		UpdatedAt: tm.UpdatedAt,
	}
}
