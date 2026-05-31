package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/pkg/utils"
)

// TenantHandler untuk handle HTTP request tenant
type TenantHandler struct {
	tenantService service.TenantService
}

// NewTenantHandler membuat instance baru TenantHandler
func NewTenantHandler(tenantService service.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// CreateTenantRequest body request untuk create tenant
type CreateTenantRequest struct {
	Name string  `json:"name" validate:"required,min=3,max=255"`
	Logo *string `json:"logo"`
}

// CreateTenant
// @Summary     Buat tenant baru
// @Description Membuat tenant baru dengan user yang login sebagai owner
// @Tags        Tenants
// @Accept      json
// @Produce     json
// @Param       body body CreateTenantRequest true "Request body"
// @Success     201 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants [post]
func (h *TenantHandler) CreateTenant(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
	}

	var req CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "name tidak boleh kosong")
	}

	// Create tenant di service
	tenant, err := h.tenantService.CreateTenant(c.Context(), userID.(string), req.Name, req.Logo)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat tenant")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, "tenant berhasil dibuat", tenant.ToResponse())
}

// GetTenants
// @Summary     List tenant user
// @Description Mengambil daftar semua tenant dimana user adalah member
// @Tags        Tenants
// @Produce     json
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants [get]
func (h *TenantHandler) GetTenants(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
	}

	tenants, err := h.tenantService.GetUserTenants(c.Context(), userID.(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil daftar tenant")
	}

	// Convert ke response
	var responses []*fiber.Map
	for _, t := range tenants {
		responses = append(responses, &fiber.Map{
			"id":         t.ID,
			"name":       t.Name,
			"slug":       t.Slug,
			"status":     t.Status,
			"logo":       t.Logo,
			"expired_at": t.ExpiredAt,
			"owner_id":   t.OwnerID,
			"created_at": t.CreatedAt,
			"updated_at": t.UpdatedAt,
		})
	}

	return utils.SuccessResponse(c, "daftar tenant berhasil diambil", responses)
}

// GetTenant
// @Summary     Detail tenant
// @Description Mengambil detail tenant berdasarkan ID (user harus member)
// @Tags        Tenants
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id} [get]
func (h *TenantHandler) GetTenant(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	tenant, err := h.tenantService.GetTenantByID(c.Context(), tenantID)
	if err != nil {
		if err.Error() == "GetTenantByID: tenant tidak ditemukan" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "tenant tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil detail tenant")
	}

	return utils.SuccessResponse(c, "detail tenant berhasil diambil", tenant.ToResponse())
}

// UpdateTenantRequest body request untuk update tenant
type UpdateTenantRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Address     *string `json:"address"`
	Description *string `json:"description"`
	Logo        *string `json:"logo"`
}

// UpdateTenant
// @Summary     Update tenant
// @Description Mengupdate data tenant (hanya admin yang bisa)
// @Tags        Tenants
// @Accept      json
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       body body UpdateTenantRequest true "Request body"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
	}

	var req UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "name tidak boleh kosong")
	}

	tenant, err := h.tenantService.UpdateTenant(c.Context(), tenantID, userID.(string), req.Name, req.Address, req.Description, req.Logo)
	if err != nil {
		if err.Error() == "UpdateTenant: user bukan admin di tenant ini" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "anda bukan admin di tenant ini")
		}
		if err.Error() == "UpdateTenant: tenant tidak ditemukan" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "tenant tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal update tenant")
	}

	return utils.SuccessResponse(c, "tenant berhasil diupdate", tenant.ToResponse())
}

// DeleteTenant
// @Summary     Delete tenant
// @Description Menghapus tenant (hanya admin yang bisa)
// @Tags        Tenants
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
	}

	err := h.tenantService.DeleteTenant(c.Context(), tenantID, userID.(string))
	if err != nil {
		if err.Error() == "DeleteTenant: user bukan admin di tenant ini" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "anda bukan admin di tenant ini")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal hapus tenant")
	}

	return utils.SuccessResponse(c, "tenant berhasil dihapus", nil)
}

// GetMembers
// @Summary     List member tenant
// @Description Mengambil daftar member di tenant (user harus member)
// @Tags        Tenants
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/members [get]
func (h *TenantHandler) GetMembers(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	members, err := h.tenantService.GetTenantMembers(c.Context(), tenantID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil daftar member")
	}

	// Convert ke response
	var responses []*fiber.Map
	for _, m := range members {
		responses = append(responses, &fiber.Map{
			"id":         m.ID,
			"tenant_id":  m.TenantID,
			"user_id":    m.UserID,
			"role":       m.Role,
			"invited_by": m.InvitedBy,
			"joined_at":  m.JoinedAt,
			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
		})
	}

	return utils.SuccessResponse(c, "daftar member berhasil diambil", responses)
}

// InviteMemberRequest body request untuk invite member
type InviteMemberRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// InviteMember
// @Summary     Invite member
// @Description Mengundang user untuk menjadi member di tenant (hanya admin yang bisa)
// @Tags        Tenants
// @Accept      json
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       body body InviteMemberRequest true "Request body"
// @Success     201 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Failure     409 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/members [post]
func (h *TenantHandler) InviteMember(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
	}

	var req InviteMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.UserID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user_id tidak boleh kosong")
	}

	member, err := h.tenantService.InviteTenantMember(c.Context(), tenantID, req.UserID, userID.(string))
	if err != nil {
		if err.Error() == "InviteTenantMember: user sudah menjadi member di tenant ini" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "user sudah menjadi member di tenant ini")
		}
		if err.Error() == "InviteTenantMember: user tidak ditemukan" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal invite member")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, "member berhasil diundang", member.ToResponse())
}

// RemoveMember
// @Summary     Remove member
// @Description Menghapus member dari tenant (hanya admin yang bisa)
// @Tags        Tenants
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       user_id path string true "User ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/members/{user_id} [delete]
func (h *TenantHandler) RemoveMember(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	memberUserID := c.Params("user_id")
	requesterID := c.Locals("user_id")
	if requesterID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan di JWT")
	}

	err := h.tenantService.RemoveTenantMember(c.Context(), tenantID, memberUserID, requesterID.(string))
	if err != nil {
		if err.Error() == "RemoveTenantMember: user bukan admin di tenant ini" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "anda bukan admin di tenant ini")
		}
		if err.Error() == "RemoveTenantMember: user bukan member di tenant ini" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user tidak menjadi member di tenant ini")
		}
		if err.Error() == "RemoveTenantMember: owner tenant tidak bisa dihapus dari member" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "owner tenant tidak bisa dihapus dari member")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal hapus member")
	}

	return utils.SuccessResponse(c, "member berhasil dihapus", nil)
}
