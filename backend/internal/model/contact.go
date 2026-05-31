package model

import (
	"time"

	"gorm.io/gorm"
)

// ContactCategory model untuk tabel contact_categories
type ContactCategory struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	TenantID  string         `gorm:"type:uuid;not null" json:"tenant_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (ContactCategory) TableName() string {
	return "contact_categories"
}

// ContactCategoryResponse untuk response API contact category
type ContactCategoryResponse struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse convert ContactCategory ke ContactCategoryResponse
func (cc *ContactCategory) ToResponse() *ContactCategoryResponse {
	return &ContactCategoryResponse{
		ID:        cc.ID,
		TenantID:  cc.TenantID,
		Name:      cc.Name,
		CreatedAt: cc.CreatedAt,
		UpdatedAt: cc.UpdatedAt,
	}
}

// Contact model untuk tabel contacts
type Contact struct {
	ID                string         `gorm:"primaryKey;type:uuid" json:"id"`
	TenantID          string         `gorm:"type:uuid;not null" json:"tenant_id"`
	ContactCategoryID string         `gorm:"type:uuid;not null" json:"contact_category_id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	IsActive          bool           `gorm:"not null;default:true" json:"is_active"`
	Address           *string        `gorm:"type:varchar(500);default:NULL" json:"address"`
	Phone             *string        `gorm:"type:varchar(50);default:NULL" json:"phone"`
	CreatedAt         time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (Contact) TableName() string {
	return "contacts"
}

// ContactResponse untuk response API contact
type ContactResponse struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenant_id"`
	ContactCategoryID string    `json:"contact_category_id"`
	Name              string    `json:"name"`
	IsActive          bool      `json:"is_active"`
	Address           *string   `json:"address"`
	Phone             *string   `json:"phone"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ToResponse convert Contact ke ContactResponse
func (c *Contact) ToResponse() *ContactResponse {
	return &ContactResponse{
		ID:                c.ID,
		TenantID:          c.TenantID,
		ContactCategoryID: c.ContactCategoryID,
		Name:              c.Name,
		IsActive:          c.IsActive,
		Address:           c.Address,
		Phone:             c.Phone,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}
