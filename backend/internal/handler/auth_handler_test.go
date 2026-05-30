package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/pkg/utils"
)

// MockAuthService mock untuk AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) RequestRegisterOTP(ctx context.Context, name, phone, password string) (string, error) {
	args := m.Called(ctx, name, phone, password)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockAuthService) VerifyRegisterOTP(ctx context.Context, phone, otp string) (*model.UserResponse, string, string, error) {
	args := m.Called(ctx, phone, otp)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*model.UserResponse), args.Get(1).(string), args.Get(2).(string), args.Error(3)
}

func (m *MockAuthService) Login(phone, password string) (*model.UserResponse, string, string, error) {
	args := m.Called(phone, password)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*model.UserResponse), args.Get(1).(string), args.Get(2).(string), args.Error(3)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, string, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(string), args.Get(1).(string), args.Error(2)
}

func (m *MockAuthService) ValidateToken(token string) (*utils.CustomClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.CustomClaims), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}

func (m *MockAuthService) GetUserByID(userID string) (*model.UserResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, phone string) (string, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, phone, otp, newPassword string) error {
	return m.Called(ctx, phone, otp, newPassword).Error(0)
}

func (m *MockAuthService) ResendRegisterOTP(ctx context.Context, phone string) (string, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(string), args.Error(1)
}

// testCfg config untuk test (non-production agar OTP tampil di response)
func testCfg() *config.AppConfig {
	return &config.AppConfig{
		App: config.AppSetting{Env: "test"},
	}
}

// Helper: setup Fiber app dengan auth handler
func setupAuthHandlerTest(authService service.AuthService) *fiber.App {
	app := fiber.New()
	authHandler := NewAuthHandler(authService, testCfg())

	// Public routes
	app.Post("/register/request-otp", authHandler.RequestRegisterOTP)
	app.Post("/register/resend-otp", authHandler.ResendRegisterOTP)
	app.Post("/register/verify", authHandler.VerifyRegisterOTP)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.RefreshToken)
	app.Post("/forgot-password", authHandler.ForgotPassword)
	app.Post("/reset-password", authHandler.ResetPassword)

	// Protected routes (untuk testing, tidak akan menggunakan middleware)
	app.Post("/logout", authHandler.Logout)
	app.Get("/me", authHandler.Me)

	return app
}

// Test RequestRegisterOTP Handler
func TestRequestRegisterOTPHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP dikirim", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("RequestRegisterOTP", mock.Anything, "John Doe", "0812345678", "password123").
			Return("123456", nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := RegisterRequest{Name: "John Doe", Phone: "0812345678", Password: "password123", ConfirmPassword: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - name kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := RegisterRequest{Name: "", Phone: "0812345678", Password: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - phone kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := RegisterRequest{Name: "John Doe", Phone: "", Password: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - password kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := RegisterRequest{Name: "John Doe", Phone: "0812345678", Password: ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - nomor sudah terdaftar", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("RequestRegisterOTP", mock.Anything, "John Doe", "0812345678", "password123").
			Return("", fmt.Errorf("RequestRegisterOTP: %w", service.ErrPhoneAlreadyRegistered))

		app := setupAuthHandlerTest(mockAuthService)

		body := RegisterRequest{Name: "John Doe", Phone: "0812345678", Password: "password123", ConfirmPassword: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - confirm_password kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := RegisterRequest{Name: "John Doe", Phone: "0812345678", Password: "password123", ConfirmPassword: ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - password tidak cocok", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := RegisterRequest{Name: "John Doe", Phone: "0812345678", Password: "password123", ConfirmPassword: "berbeda456"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - body tidak valid", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		req := httptest.NewRequest(http.MethodPost, "/register/request-otp", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

// Test VerifyRegisterOTP Handler
func TestVerifyRegisterOTPHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - registrasi berhasil", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		userResp := &model.UserResponse{ID: "user-123", Name: "John Doe", Phone: "6212345678"}
		mockAuthService.On("VerifyRegisterOTP", mock.Anything, "0812345678", "123456").
			Return(userResp, "access-token", "refresh-token", nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := VerifyOTPRequest{Phone: "0812345678", OTP: "123456"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - phone kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := VerifyOTPRequest{Phone: "", OTP: "123456"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - OTP kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := VerifyOTPRequest{Phone: "0812345678", OTP: ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - OTP salah", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("VerifyRegisterOTP", mock.Anything, "0812345678", "999999").
			Return(nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", service.ErrOTPInvalid))

		app := setupAuthHandlerTest(mockAuthService)

		body := VerifyOTPRequest{Phone: "0812345678", OTP: "999999"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - OTP expired", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("VerifyRegisterOTP", mock.Anything, "0812345678", "123456").
			Return(nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", service.ErrOTPExpired))

		app := setupAuthHandlerTest(mockAuthService)

		body := VerifyOTPRequest{Phone: "0812345678", OTP: "123456"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - data pending expired", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("VerifyRegisterOTP", mock.Anything, "0812345678", "123456").
			Return(nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", service.ErrPendingDataExpired))

		app := setupAuthHandlerTest(mockAuthService)

		body := VerifyOTPRequest{Phone: "0812345678", OTP: "123456"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

// Test Login Handler
func TestLoginHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - login berhasil", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		userResp := &model.UserResponse{ID: "user-123", Name: "John Doe", Phone: "6212345678"}
		mockAuthService.On("Login", "0812345678", "password123").
			Return(userResp, "access-token", "refresh-token", nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := LoginRequest{Phone: "0812345678", Password: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - phone kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := LoginRequest{Phone: "", Password: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - password kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := LoginRequest{Phone: "0812345678", Password: ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - user tidak ditemukan", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("Login", "0812345678", "password123").
			Return(nil, "", "", fmt.Errorf("Login: %w", service.ErrUserNotFound))

		app := setupAuthHandlerTest(mockAuthService)

		body := LoginRequest{Phone: "0812345678", Password: "password123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - password salah", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("Login", "0812345678", "wrongpassword").
			Return(nil, "", "", fmt.Errorf("Login: %w", service.ErrPasswordIncorrect))

		app := setupAuthHandlerTest(mockAuthService)

		body := LoginRequest{Phone: "0812345678", Password: "wrongpassword"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

// Test RefreshToken Handler
func TestRefreshTokenHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - refresh token berhasil", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("RefreshToken", "old-refresh-token").
			Return("new-access-token", "new-refresh-token", nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := RefreshTokenRequest{RefreshToken: "old-refresh-token"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - refresh_token kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := RefreshTokenRequest{RefreshToken: ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - token di-blacklist", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("RefreshToken", "blacklisted-token").
			Return("", "", fmt.Errorf("RefreshToken: %w", service.ErrTokenBlacklisted))

		app := setupAuthHandlerTest(mockAuthService)

		body := RefreshTokenRequest{RefreshToken: "blacklisted-token"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

// Test ForgotPassword Handler
func TestForgotPasswordHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP dikirim", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ForgotPassword", mock.Anything, "0812345678").
			Return("123456", nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := ForgotPasswordRequest{Phone: "0812345678"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - phone kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := ForgotPasswordRequest{Phone: ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - phone tidak terdaftar", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ForgotPassword", mock.Anything, "0812345678").
			Return("", fmt.Errorf("ForgotPassword: %w", service.ErrUserNotFound))

		app := setupAuthHandlerTest(mockAuthService)

		body := ForgotPasswordRequest{Phone: "0812345678"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

// Test ResetPassword Handler
func TestResetPasswordHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - password direset", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ResetPassword", mock.Anything, "0812345678", "123456", "newpassword123").
			Return(nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := ResetPasswordRequest{Phone: "0812345678", OTP: "123456", NewPassword: "newpassword123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - phone kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := ResetPasswordRequest{Phone: "", OTP: "123456", NewPassword: "newpassword123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - OTP tidak valid", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ResetPassword", mock.Anything, "0812345678", "999999", "newpassword123").
			Return(fmt.Errorf("ResetPassword: %w", service.ErrOTPInvalid))

		app := setupAuthHandlerTest(mockAuthService)

		body := ResetPasswordRequest{Phone: "0812345678", OTP: "999999", NewPassword: "newpassword123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - OTP expired", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ResetPassword", mock.Anything, "0812345678", "123456", "newpassword123").
			Return(fmt.Errorf("ResetPassword: %w", service.ErrOTPExpired))

		app := setupAuthHandlerTest(mockAuthService)

		body := ResetPasswordRequest{Phone: "0812345678", OTP: "123456", NewPassword: "newpassword123"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

// Test Logout Handler
func TestLogoutHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - logout berhasil", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("Logout", mock.Anything, "valid-token").Return(nil)

		app := setupAuthHandlerTest(mockAuthService)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - authorization header tidak ada", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

// Test Me Handler
func TestMeHandler(t *testing.T) {
	t.Parallel()

	t.Run("Gagal - user_id tidak ada di Locals", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		req := httptest.NewRequest(http.MethodGet, "/me", nil)

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

// Test ResendRegisterOTP Handler
func TestResendRegisterOTPHandler(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - OTP berhasil dikirim ulang", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ResendRegisterOTP", mock.Anything, "0812345678").
			Return("654321", nil)

		app := setupAuthHandlerTest(mockAuthService)

		body := map[string]string{"phone": "0812345678"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/resend-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - phone kosong", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		body := map[string]string{"phone": ""}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/resend-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Gagal - tidak ada registrasi pending", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)
		mockAuthService.On("ResendRegisterOTP", mock.Anything, "0812345678").
			Return("", fmt.Errorf("ResendRegisterOTP: %w", service.ErrNoPendingRegistration))

		app := setupAuthHandlerTest(mockAuthService)

		body := map[string]string{"phone": "0812345678"}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/register/resend-otp", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Gagal - body tidak valid", func(t *testing.T) {
		t.Parallel()

		app := setupAuthHandlerTest(new(MockAuthService))

		req := httptest.NewRequest(http.MethodPost, "/register/resend-otp", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
