package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository mock untuk UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByPhone(phone string) (*model.User, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(userID string) error {
	return m.Called(userID).Error(0)
}

func (m *MockUserRepository) UpdatePassword(userID, hashedPassword string) error {
	return m.Called(userID, hashedPassword).Error(0)
}

func (m *MockUserRepository) UpdateVerifiedAt(userID string) error {
	return m.Called(userID).Error(0)
}

func (m *MockUserRepository) GetAll() ([]model.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(userID string) error {
	return m.Called(userID).Error(0)
}

// MockNotificationService mock untuk NotificationService interface
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) Send(ctx context.Context, phone, message string) error {
	return m.Called(ctx, phone, message).Error(0)
}

// setupTest helper untuk membuat test environment
func setupTest(t *testing.T) (*miniredis.Miniredis, *redis.Client, *config.AppConfig) {
	// Setup miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)

	// Setup redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Setup config
	cfg := &config.AppConfig{
		App: config.AppSetting{
			Env: "test",
		},
		JWT: config.JWTConfig{
			Secret:           "test-secret-key-very-secure-123",
			ExpiryHour:       1,
			RefreshExpiryDay: 7,
		},
	}

	return mr, rdb, cfg
}

// Test RequestRegisterOTP
func TestRequestRegisterOTPAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP digenerate dan notifikasi dikirim", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		mockUserRepo.On("FindByPhone", "6212345678").Return(nil, nil)

		// Verifikasi Send dipanggil dengan phone yang dinormalisasi dan pesan berisi "Kode OTP"
		mockNotifService.On("Send", mock.Anything, "6212345678", mock.MatchedBy(func(msg string) bool {
			return len(msg) > 0 && msg[0:4] == "Kode"
		})).Return(nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.RequestRegisterOTP(ctx, "John Doe", "0812345678", "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		assert.Len(t, otp, 6)
		mockUserRepo.AssertExpectations(t)
		mockNotifService.AssertExpectations(t)
	})

	t.Run("Sukses - notifikasi gagal tapi OTP tetap disimpan", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		mockUserRepo.On("FindByPhone", "6212345678").Return(nil, nil)
		mockNotifService.On("Send", mock.Anything, mock.Anything, mock.Anything).
			Return(fmt.Errorf("WhatsApp API error"))

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.RequestRegisterOTP(ctx, "John Doe", "0812345678", "password123")

		// Flow tetap sukses meskipun notifikasi gagal
		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		mockNotifService.AssertExpectations(t)
	})

	t.Run("Sukses - notificationService nil tidak panic", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockUserRepo.On("FindByPhone", "6212345678").Return(nil, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, nil, nil)

		ctx := context.Background()
		otp, err := authService.RequestRegisterOTP(ctx, "John Doe", "0812345678", "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, otp)
	})

	t.Run("Gagal - phone sudah terdaftar", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		existingUser := &model.User{
			ID:    "user-123",
			Name:  "Existing User",
			Phone: "6212345678",
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(existingUser, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.RequestRegisterOTP(ctx, "John Doe", "0812345678", "password123")

		assert.Error(t, err)
		assert.Empty(t, otp)
		mockUserRepo.AssertExpectations(t)
		// Send tidak boleh dipanggil
		mockNotifService.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything)
	})
}

// Test VerifyRegisterOTP
func TestVerifyRegisterOTPAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP valid, user dibuat", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		// Setup data di Redis seperti hasil RequestRegisterOTP
		ctx := context.Background()
		otpKey := "otp:register:6212345678"
		pendingKey := "register:pending:6212345678"
		hashedOTP, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		mr.Set(otpKey, string(hashedOTP))
		mr.Set(pendingKey, `{"name":"John Doe","password":"$2a$10$hashedpassword123"}`)

		mockUserRepo.On("Create", mock.MatchedBy(func(u *model.User) bool {
			return u.Name == "John Doe" && u.Phone == "6212345678" && u.IsAdmin == false && u.VerifiedAt != nil
		})).Return(nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, accessToken, refreshToken, err := authService.VerifyRegisterOTP(ctx, "0812345678", "123456")

		require.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.Equal(t, "John Doe", userResp.Name)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Gagal - OTP tidak valid", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		ctx := context.Background()
		otpKey := "otp:register:6212345678"
		mr.Set(otpKey, "123456")

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, accessToken, refreshToken, err := authService.VerifyRegisterOTP(ctx, "0812345678", "999999")

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})

	t.Run("Gagal - OTP tidak ditemukan / expired", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		userResp, accessToken, refreshToken, err := authService.VerifyRegisterOTP(ctx, "0812345678", "123456")

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})

	t.Run("Gagal - data pending tidak ditemukan", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		ctx := context.Background()
		otpKey := "otp:register:6212345678"
		mr.Set(otpKey, "123456")
		// Tidak set pendingKey sehingga data pending tidak ada

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, accessToken, refreshToken, err := authService.VerifyRegisterOTP(ctx, "0812345678", "123456")

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})
}

// Test Login
func TestLoginAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - login dengan phone + password benar", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		hashedStr := string(hashedPassword)

		user := &model.User{
			ID:       "user-123",
			Name:     "John Doe",
			Phone:    "6212345678",
			IsAdmin:  false,
			Password: &hashedStr,
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)
		mockUserRepo.On("UpdateLastLogin", "user-123").Return(nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, accessToken, refreshToken, err := authService.Login("0812345678", "password123")

		require.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Gagal - user tidak ditemukan", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		mockUserRepo.On("FindByPhone", "6212345678").Return(nil, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, accessToken, refreshToken, err := authService.Login("0812345678", "password123")

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})

	t.Run("Gagal - password salah", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		hashedStr := string(hashedPassword)

		user := &model.User{
			ID:       "user-123",
			Name:     "John Doe",
			Phone:    "6212345678",
			IsAdmin:  false,
			Password: &hashedStr,
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, accessToken, refreshToken, err := authService.Login("0812345678", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})

	t.Run("Gagal - user tidak punya password", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		user := &model.User{
			ID:       "user-123",
			Name:     "John Doe",
			Phone:    "6212345678",
			IsAdmin:  false,
			Password: nil,
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, _, _, err := authService.Login("0812345678", "password123")

		assert.Error(t, err)
		assert.Nil(t, userResp)
	})
}

// Test RefreshToken
func TestRefreshTokenAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - refresh token valid", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		user := &model.User{
			ID:      "user-123",
			Name:    "John Doe",
			IsAdmin: false,
		}

		// Generate token pair
		_, refreshToken, err := utils.GenerateTokenPair("user-123", false, cfg.JWT.Secret, cfg.JWT.ExpiryHour, cfg.JWT.RefreshExpiryDay)
		require.NoError(t, err)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)
		mockUserRepo.On("FindByID", "user-123").Return(user, nil)

		newAccessToken, newRefreshToken, err := authService.RefreshToken(refreshToken)

		require.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Gagal - token bukan refresh token", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		// Generate access token (bukan refresh token)
		accessToken, _, _ := utils.GenerateTokenPair("user-123", false, cfg.JWT.Secret, cfg.JWT.ExpiryHour, cfg.JWT.RefreshExpiryDay)

		_, _, err := authService.RefreshToken(accessToken)

		assert.Error(t, err)
	})
}

// Test ValidateToken
func TestValidateTokenAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - validate access token", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		accessToken, _, _ := utils.GenerateTokenPair("user-123", false, cfg.JWT.Secret, cfg.JWT.ExpiryHour, cfg.JWT.RefreshExpiryDay)

		claims, err := authService.ValidateToken(accessToken)

		require.NoError(t, err)
		assert.Equal(t, "user-123", claims.UserID)
		assert.False(t, claims.IsAdmin)
		assert.Equal(t, "access", claims.Type)
	})

	t.Run("Gagal - validate refresh token sebagai access token", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		_, refreshToken, _ := utils.GenerateTokenPair("user-123", false, cfg.JWT.Secret, cfg.JWT.ExpiryHour, cfg.JWT.RefreshExpiryDay)

		_, err := authService.ValidateToken(refreshToken)

		assert.Error(t, err)
	})
}

// Test Logout
func TestLogoutAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - logout menambahkan token ke blacklist", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		accessToken, _, _ := utils.GenerateTokenPair("user-123", false, cfg.JWT.Secret, cfg.JWT.ExpiryHour, cfg.JWT.RefreshExpiryDay)

		ctx := context.Background()
		err := authService.Logout(ctx, accessToken)

		require.NoError(t, err)

		// Verify token is blacklisted
		_, err = authService.ValidateToken(accessToken)
		assert.Error(t, err)
	})
}

// Test GetUserByID
func TestGetUserByIDAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - dapatkan user berdasarkan ID", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		user := &model.User{
			ID:      "user-123",
			Name:    "John Doe",
			Phone:   "6212345678",
			IsAdmin: false,
		}

		mockUserRepo.On("FindByID", "user-123").Return(user, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, err := authService.GetUserByID("user-123")

		require.NoError(t, err)
		assert.Equal(t, "user-123", userResp.ID)
		assert.Equal(t, "John Doe", userResp.Name)
	})

	t.Run("Gagal - user tidak ditemukan", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		mockUserRepo.On("FindByID", "user-123").Return(nil, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		userResp, err := authService.GetUserByID("user-123")

		assert.Error(t, err)
		assert.Nil(t, userResp)
	})
}

// Test ForgotPassword
func TestForgotPasswordAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP digenerate dan notifikasi dikirim", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		user := &model.User{
			ID:    "user-123",
			Name:  "John Doe",
			Phone: "6212345678",
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)
		// Verifikasi Send dipanggil dengan phone yang dinormalisasi dan pesan berisi "Kode OTP"
		mockNotifService.On("Send", mock.Anything, "6212345678", mock.MatchedBy(func(msg string) bool {
			return len(msg) > 0 && msg[0:4] == "Kode"
		})).Return(nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.ForgotPassword(ctx, "0812345678")

		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		assert.Len(t, otp, 6)
		mockUserRepo.AssertExpectations(t)
		mockNotifService.AssertExpectations(t)
	})

	t.Run("Sukses - notifikasi gagal tapi OTP tetap disimpan", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		user := &model.User{
			ID:    "user-123",
			Name:  "John Doe",
			Phone: "6212345678",
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)
		mockNotifService.On("Send", mock.Anything, mock.Anything, mock.Anything).
			Return(fmt.Errorf("WhatsApp API error"))

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.ForgotPassword(ctx, "0812345678")

		// Flow tetap sukses meskipun notifikasi gagal
		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		mockNotifService.AssertExpectations(t)
	})

	t.Run("Sukses - notificationService nil tidak panic", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)

		user := &model.User{
			ID:    "user-123",
			Name:  "John Doe",
			Phone: "6212345678",
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, nil, nil)

		ctx := context.Background()
		otp, err := authService.ForgotPassword(ctx, "0812345678")

		require.NoError(t, err)
		assert.NotEmpty(t, otp)
	})

	t.Run("Gagal - phone tidak terdaftar, Send tidak dipanggil", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		mockUserRepo.On("FindByPhone", "6212345678").Return(nil, nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.ForgotPassword(ctx, "0812345678")

		assert.Error(t, err)
		assert.Empty(t, otp)
		mockUserRepo.AssertExpectations(t)
		// Send tidak boleh dipanggil
		mockNotifService.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything)
	})
}

// Test ResetPassword
func TestResetPasswordAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - reset password dengan OTP valid", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		// Setup OTP di Redis
		ctx := context.Background()
		otpKey := "otp:reset:6212345678"
		hashedOTP, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		mr.Set(otpKey, string(hashedOTP))

		user := &model.User{
			ID:    "user-123",
			Name:  "John Doe",
			Phone: "6212345678",
		}

		mockUserRepo.On("FindByPhone", "6212345678").Return(user, nil)
		mockUserRepo.On("UpdatePassword", "user-123", mock.AnythingOfType("string")).Return(nil)

		err := authService.ResetPassword(ctx, "0812345678", "123456", "newpassword123")

		require.NoError(t, err)
	})

	t.Run("Gagal - OTP tidak valid", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		// Setup OTP di Redis
		ctx := context.Background()
		otpKey := "otp:reset:6212345678"
		hashedOTP, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		mr.Set(otpKey, string(hashedOTP))

		err := authService.ResetPassword(ctx, "0812345678", "999999", "newpassword123")

		assert.Error(t, err)
	})

	t.Run("Gagal - OTP tidak ditemukan / expired", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		err := authService.ResetPassword(ctx, "0812345678", "123456", "newpassword123")

		assert.Error(t, err)
	})
}

// Test ResendRegisterOTP
func TestResendRegisterOTPAuthService(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP baru digenerate dan dikirim", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		ctx := context.Background()
		// Simulasikan data pending yang ada
		mr.Set("register:pending:6212345678", `{"name":"John Doe","password":"$2a$10$hashedpassword"}`)

		mockNotifService.On("Send", mock.Anything, "6212345678", mock.MatchedBy(func(msg string) bool {
			return len(msg) > 0 && msg[:4] == "Kode"
		})).Return(nil)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		otp, err := authService.ResendRegisterOTP(ctx, "0812345678")

		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		assert.Len(t, otp, 6)
		mockNotifService.AssertExpectations(t)
	})

	t.Run("Sukses - notifikasi gagal tapi OTP tetap disimpan", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		ctx := context.Background()
		mr.Set("register:pending:6212345678", `{"name":"John Doe","password":"$2a$10$hashedpassword"}`)

		mockNotifService.On("Send", mock.Anything, mock.Anything, mock.Anything).
			Return(fmt.Errorf("WhatsApp API error"))

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		otp, err := authService.ResendRegisterOTP(ctx, "0812345678")

		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		mockNotifService.AssertExpectations(t)
	})

	t.Run("Gagal - tidak ada data pending", func(t *testing.T) {
		t.Parallel()

		mr, rdb, cfg := setupTest(t)
		defer mr.Close()
		defer rdb.Close()

		mockUserRepo := new(MockUserRepository)
		mockNotifService := new(MockNotificationService)

		authService := NewAuthService(mockUserRepo, rdb, cfg, mockNotifService, nil)

		ctx := context.Background()
		otp, err := authService.ResendRegisterOTP(ctx, "0812345678")

		assert.Error(t, err)
		assert.Empty(t, otp)
		assert.ErrorIs(t, err, ErrNoPendingRegistration)
		mockNotifService.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything)
	})
}
