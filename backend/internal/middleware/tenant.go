package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/username/qurban-app/internal/repository"
	"github.com/username/qurban-app/pkg/utils"
)

// TenantMember middleware untuk validasi apakah user adalah member di tenant
func TenantMember(tenantRepo repository.TenantRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil tenant_id dari URL param
		tenantID := c.Params("id")
		if tenantID == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "tenant id tidak ditemukan")
		}

		// Ambil user_id dari JWT locals (harus sudah di-set oleh JWTAuth middleware)
		userID := c.Locals("user_id")
		if userID == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
		}

		// Cek apakah user adalah member di tenant
		isMember, err := tenantRepo.IsTenantMember(c.Context(), tenantID, userID.(string))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal validasi membership: "+err.Error())
		}

		if !isMember {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "anda bukan member di tenant ini")
		}

		// Set tenant_id ke locals untuk digunakan di handler
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

// TenantAdmin middleware untuk validasi apakah user adalah admin di tenant
func TenantAdmin(tenantRepo repository.TenantRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil tenant_id dari URL param
		tenantID := c.Params("id")
		if tenantID == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "tenant id tidak ditemukan")
		}

		// Ambil user_id dari JWT locals (harus sudah di-set oleh JWTAuth middleware)
		userID := c.Locals("user_id")
		if userID == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
		}

		// Cek apakah user adalah admin di tenant
		isAdmin, err := tenantRepo.IsTenantAdmin(c.Context(), tenantID, userID.(string))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal validasi admin status: "+err.Error())
		}

		if !isAdmin {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "anda bukan admin di tenant ini")
		}

		// Set tenant_id ke locals untuk digunakan di handler
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}
