package model

import (
	"time"

	"gorm.io/gorm"
)

// User model untuk tabel users
type User struct {
	ID          string             `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string             `gorm:"type:varchar(255);not null" json:"name"`
	Phone       string             `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Email       *string            `gorm:"type:varchar(255);uniqueIndex;default:NULL" json:"email"`
	IsAdmin     bool               `gorm:"type:boolean;default:false" json:"is_admin"`
	Password    *string            `gorm:"type:varchar(255);default:NULL" json:"-"`
	VerifiedAt  *time.Time         `gorm:"type:timestamp;default:NULL" json:"verified_at"`
	LastLoginAt *time.Time         `gorm:"type:timestamp;default:NULL" json:"last_login_at"`
	CreatedAt   time.Time          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (User) TableName() string {
	return "users"
}

// UserResponse untuk response API (tanpa password)
type UserResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	Email       *string    `json:"email"`
	IsAdmin     bool       `json:"is_admin"`
	VerifiedAt  *time.Time `json:"verified_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse convert User ke UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Phone:       u.Phone,
		Email:       u.Email,
		IsAdmin:     u.IsAdmin,
		VerifiedAt:  u.VerifiedAt,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
