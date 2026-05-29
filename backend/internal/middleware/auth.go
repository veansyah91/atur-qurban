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
		// Ambil Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "authorization header tidak ditemukan")
		}

		// Extract token (remove "Bearer " prefix)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "format authorization header salah")
		}

		token := parts[1]

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
