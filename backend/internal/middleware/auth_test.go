package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/pkg/utils"
)

// MockAuthService mock untuk AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) RequestRegisterOTP(ctx context.Context, name, phone, password string) (string, error) {
	args := m.Called(ctx, name, phone, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ResendRegisterOTP(ctx context.Context, phone string) (string, error) {
	args := m.Called(ctx, phone)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyRegisterOTP(ctx context.Context, phone, otp string) (*model.UserResponse, string, string, error) {
	args := m.Called(ctx, phone, otp)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*model.UserResponse), args.String(1), args.String(2), args.Error(3)
}

func (m *MockAuthService) Login(phone, password string) (*model.UserResponse, string, string, error) {
	args := m.Called(phone, password)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*model.UserResponse), args.String(1), args.String(2), args.Error(3)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.String(1), args.Error(2)
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
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, phone, otp, newPassword string) error {
	return m.Called(ctx, phone, otp, newPassword).Error(0)
}

// Test JWTAuth Middleware
func TestJWTAuth(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - valid token dengan Authorization header", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		claims := &utils.CustomClaims{
			UserID:  "user-123",
			IsAdmin: false,
			Type:    "access",
		}

		mockAuthService.On("ValidateToken", "valid-token").Return(claims, nil)

		app := fiber.New()
		app.Use(JWTAuth(mockAuthService))

		app.Get("/protected", func(c *fiber.Ctx) error {
			userID := c.Locals("user_id")
			isAdmin := c.Locals("is_admin")

			return c.JSON(fiber.Map{
				"user_id": userID,
				"is_admin": isAdmin,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertCalled(t, "ValidateToken", "valid-token")
	})

	t.Run("Gagal - Authorization header tidak ada", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		app := fiber.New()
		app.Use(JWTAuth(mockAuthService))

		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Gagal - Authorization header format salah", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		app := fiber.New()
		app.Use(JWTAuth(mockAuthService))

		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat valid-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Gagal - token tidak valid", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		mockAuthService.On("ValidateToken", "invalid-token").Return(nil, fiber.NewError(fiber.StatusUnauthorized, "token tidak valid"))

		app := fiber.New()
		app.Use(JWTAuth(mockAuthService))

		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		mockAuthService.AssertCalled(t, "ValidateToken", "invalid-token")
	})

	t.Run("Sukses - admin user dapat mengakses protected route", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		adminClaims := &utils.CustomClaims{
			UserID:  "admin-user",
			IsAdmin: true,
			Type:    "access",
		}

		mockAuthService.On("ValidateToken", "admin-token").Return(adminClaims, nil)

		app := fiber.New()
		app.Use(JWTAuth(mockAuthService))

		app.Get("/protected", func(c *fiber.Ctx) error {
			isAdmin := c.Locals("is_admin")
			return c.JSON(fiber.Map{
				"is_admin": isAdmin,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer admin-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

// Test AdminOnly Middleware
func TestAdminOnly(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - admin dapat mengakses", func(t *testing.T) {
		t.Parallel()

		app := fiber.New()

		app.Get("/admin", func(c *fiber.Ctx) error {
			// Simulasi adalah_admin = true di Locals (set oleh JWTAuth middleware)
			c.Locals("is_admin", true)
			return c.Next()
		}, AdminOnly, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Gagal - non-admin tidak dapat mengakses", func(t *testing.T) {
		t.Parallel()

		app := fiber.New()

		app.Get("/admin", func(c *fiber.Ctx) error {
			// Simulasi is_admin = false
			c.Locals("is_admin", false)
			return c.Next()
		}, AdminOnly, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Gagal - is_admin tidak ada di Locals", func(t *testing.T) {
		t.Parallel()

		app := fiber.New()

		app.Get("/admin", AdminOnly, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Gagal - is_admin adalah nil", func(t *testing.T) {
		t.Parallel()

		app := fiber.New()

		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("is_admin", nil)
			return c.Next()
		}, AdminOnly, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})
}

// Integration test: JWTAuth + AdminOnly
func TestJWTAuthWithAdminOnly(t *testing.T) {
	t.Parallel()

	t.Run("Admin dapat mengakses admin route dengan valid token", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		adminClaims := &utils.CustomClaims{
			UserID:  "admin-user",
			IsAdmin: true,
			Type:    "access",
		}

		mockAuthService.On("ValidateToken", "admin-token").Return(adminClaims, nil)

		app := fiber.New()

		app.Get("/admin", JWTAuth(mockAuthService), AdminOnly, func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "admin access granted",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer admin-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthService.AssertCalled(t, "ValidateToken", "admin-token")
	})

	t.Run("User biasa tidak dapat mengakses admin route", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		userClaims := &utils.CustomClaims{
			UserID:  "user-123",
			IsAdmin: false,
			Type:    "access",
		}

		mockAuthService.On("ValidateToken", "user-token").Return(userClaims, nil)

		app := fiber.New()

		app.Get("/admin", JWTAuth(mockAuthService), AdminOnly, func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "admin access granted",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer user-token")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
		mockAuthService.AssertCalled(t, "ValidateToken", "user-token")
	})

	t.Run("Request tanpa token tidak dapat mengakses admin route", func(t *testing.T) {
		t.Parallel()

		mockAuthService := new(MockAuthService)

		app := fiber.New()

		app.Get("/admin", JWTAuth(mockAuthService), AdminOnly, func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "admin access granted",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
