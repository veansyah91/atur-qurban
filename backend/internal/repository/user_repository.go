package repository

import (
	"fmt"

	"github.com/username/qurban-app/internal/model"
	"gorm.io/gorm"
)

// UserRepository interface untuk operasi user
type UserRepository interface {
	FindByPhone(phone string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Create(user *model.User) error
	UpdateLastLogin(userID string) error
	UpdatePassword(userID, hashedPassword string) error
	GetAll() ([]model.User, error)
	Delete(userID string) error
}

// userRepository implementasi UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByPhone mencari user berdasarkan nomor telepon
func (r *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User tidak ditemukan, return nil (bukan error)
		}
		return nil, fmt.Errorf("FindByPhone: %w", err)
	}
	return &user, nil
}

// FindByID mencari user berdasarkan ID
func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return &user, nil
}

// Create menyimpan user baru ke database
func (r *userRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("Create: %w", err)
	}
	return nil
}

// UpdateLastLogin mengupdate waktu login terakhir user
func (r *userRepository) UpdateLastLogin(userID string) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return fmt.Errorf("UpdateLastLogin: %w", err)
	}
	return nil
}

// GetAll mengambil semua user
func (r *userRepository) GetAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("GetAll: %w", err)
	}
	return users, nil
}

// Delete menghapus user secara soft delete
func (r *userRepository) Delete(userID string) error {
	if err := r.db.Delete(&model.User{}, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}

// UpdatePassword mengupdate password user
func (r *userRepository) UpdatePassword(userID, hashedPassword string) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("UpdatePassword: %w", err)
	}
	return nil
}
