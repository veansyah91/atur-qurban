package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/pkg/utils"
)

// JWTAuth middleware untuk validasi JWT token
func JWTAuth(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari cookie atau Authorization header
		token := c.Cookies("access_token")

		if token == "" {
			// Fallback ke Authorization header
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		if token == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "autentikasi diperlukan")
		}

		// Validasi token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "token tidak valid: "+err.Error())
		}

		// Set locals untuk digunakan di handler
		c.Locals("user_id", claims.UserID)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}

// AdminOnly middleware untuk validasi admin
func AdminOnly(c *fiber.Ctx) error {
	isAdmin := c.Locals("is_admin")
	if isAdmin == nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "akses ditolak: anda bukan admin")
	}

	if !isAdmin.(bool) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "akses ditolak: anda bukan admin")
	}

	return c.Next()
}

// RequireVerified middleware untuk memastikan user sudah terverifikasi
func RequireVerified(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan")
		}

		// Ambil data user dari service
		userResp, err := authService.GetUserByID(userID.(string))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "gagal ambil data user: "+err.Error())
		}

		// Cek apakah user sudah terverifikasi
		if userResp.VerifiedAt == nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "akun anda belum terverifikasi")
		}

		return c.Next()
	}
}
