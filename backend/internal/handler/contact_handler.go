package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/pkg/utils"
)

// ContactCategoryHandler untuk handle HTTP request contact category
type ContactCategoryHandler struct {
	categoryService service.ContactCategoryService
}

// NewContactCategoryHandler membuat instance baru ContactCategoryHandler
func NewContactCategoryHandler(categoryService service.ContactCategoryService) *ContactCategoryHandler {
	return &ContactCategoryHandler{categoryService: categoryService}
}

// CreateContactCategoryRequest body request untuk create contact category
type CreateContactCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

// UpdateContactCategoryRequest body request untuk update contact category
type UpdateContactCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

// ListCategories
// @Summary     List kategori kontak
// @Description Mengambil daftar semua kategori kontak dalam tenant
// @Tags        ContactCategories
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contact-categories [get]
func (h *ContactCategoryHandler) ListCategories(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	categories, err := h.categoryService.GetByTenant(c.Context(), tenantID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil daftar kategori kontak")
	}

	var responses []interface{}
	for _, cat := range categories {
		responses = append(responses, cat.ToResponse())
	}
	if responses == nil {
		responses = []interface{}{}
	}

	return utils.SuccessResponse(c, "daftar kategori kontak berhasil diambil", responses)
}

// CreateCategory
// @Summary     Buat kategori kontak
// @Description Membuat kategori kontak baru dalam tenant (hanya admin)
// @Tags        ContactCategories
// @Accept      json
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       body body CreateContactCategoryRequest true "Request body"
// @Success     201 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contact-categories [post]
func (h *ContactCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	var req CreateContactCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "name tidak boleh kosong")
	}

	category, err := h.categoryService.Create(c.Context(), tenantID, req.Name)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat kategori kontak")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, "kategori kontak berhasil dibuat", category.ToResponse())
}

// UpdateCategory
// @Summary     Update kategori kontak
// @Description Mengupdate nama kategori kontak (hanya admin)
// @Tags        ContactCategories
// @Accept      json
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       cat_id path string true "Category ID"
// @Param       body body UpdateContactCategoryRequest true "Request body"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contact-categories/{cat_id} [put]
func (h *ContactCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	catID := c.Params("cat_id")

	var req UpdateContactCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "name tidak boleh kosong")
	}

	category, err := h.categoryService.Update(c.Context(), catID, tenantID, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kategori kontak tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal update kategori kontak")
	}

	return utils.SuccessResponse(c, "kategori kontak berhasil diupdate", category.ToResponse())
}

// DeleteCategory
// @Summary     Hapus kategori kontak
// @Description Menghapus kategori kontak secara soft delete (hanya admin)
// @Tags        ContactCategories
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       cat_id path string true "Category ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contact-categories/{cat_id} [delete]
func (h *ContactCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	catID := c.Params("cat_id")

	err := h.categoryService.Delete(c.Context(), catID, tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kategori kontak tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal hapus kategori kontak")
	}

	return utils.SuccessResponse(c, "kategori kontak berhasil dihapus", nil)
}

// ─────────────────────────────────────────────
// ContactHandler
// ─────────────────────────────────────────────

// ContactHandler untuk handle HTTP request contact
type ContactHandler struct {
	contactService service.ContactService
}

// NewContactHandler membuat instance baru ContactHandler
func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

// CreateContactRequest body request untuk create contact
type CreateContactRequest struct {
	ContactCategoryID string  `json:"contact_category_id" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	IsActive          *bool   `json:"is_active"`
	Address           *string `json:"address"`
	Phone             *string `json:"phone"`
}

// UpdateContactRequest body request untuk update contact
type UpdateContactRequest struct {
	ContactCategoryID string  `json:"contact_category_id" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	IsActive          bool    `json:"is_active"`
	Address           *string `json:"address"`
	Phone             *string `json:"phone"`
}

// ListContacts
// @Summary     List kontak
// @Description Mengambil daftar semua kontak dalam tenant
// @Tags        Contacts
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contacts [get]
func (h *ContactHandler) ListContacts(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	contacts, err := h.contactService.GetByTenant(c.Context(), tenantID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil daftar kontak")
	}

	var responses []interface{}
	for _, con := range contacts {
		responses = append(responses, con.ToResponse())
	}
	if responses == nil {
		responses = []interface{}{}
	}

	return utils.SuccessResponse(c, "daftar kontak berhasil diambil", responses)
}

// CreateContact
// @Summary     Buat kontak baru
// @Description Membuat kontak baru dalam tenant (hanya admin)
// @Tags        Contacts
// @Accept      json
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       body body CreateContactRequest true "Request body"
// @Success     201 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contacts [post]
func (h *ContactHandler) CreateContact(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	var req CreateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "name tidak boleh kosong")
	}
	if req.ContactCategoryID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "contact_category_id tidak boleh kosong")
	}

	contact, err := h.contactService.Create(c.Context(), tenantID, req.ContactCategoryID, req.Name, req.Address, req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kategori kontak tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal membuat kontak")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, "kontak berhasil dibuat", contact.ToResponse())
}

// GetContact
// @Summary     Detail kontak
// @Description Mengambil detail kontak berdasarkan ID
// @Tags        Contacts
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       contact_id path string true "Contact ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contacts/{contact_id} [get]
func (h *ContactHandler) GetContact(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	contactID := c.Params("contact_id")

	contact, err := h.contactService.GetByID(c.Context(), contactID, tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kontak tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil detail kontak")
	}

	return utils.SuccessResponse(c, "detail kontak berhasil diambil", contact.ToResponse())
}

// UpdateContact
// @Summary     Update kontak
// @Description Mengupdate data kontak (hanya admin)
// @Tags        Contacts
// @Accept      json
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       contact_id path string true "Contact ID"
// @Param       body body UpdateContactRequest true "Request body"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contacts/{contact_id} [put]
func (h *ContactHandler) UpdateContact(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	contactID := c.Params("contact_id")

	var req UpdateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "format body tidak valid")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "name tidak boleh kosong")
	}
	if req.ContactCategoryID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "contact_category_id tidak boleh kosong")
	}

	contact, err := h.contactService.Update(c.Context(), contactID, tenantID, req.ContactCategoryID, req.Name, req.IsActive, req.Address, req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kontak tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal update kontak")
	}

	return utils.SuccessResponse(c, "kontak berhasil diupdate", contact.ToResponse())
}

// DeleteContact
// @Summary     Hapus kontak
// @Description Menghapus kontak secara soft delete (hanya admin)
// @Tags        Contacts
// @Produce     json
// @Param       id path string true "Tenant ID"
// @Param       contact_id path string true "Contact ID"
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     403 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Security    BearerAuth
// @Router      /api/v1/tenants/{id}/contacts/{contact_id} [delete]
func (h *ContactHandler) DeleteContact(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	contactID := c.Params("contact_id")

	err := h.contactService.Delete(c.Context(), contactID, tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "kontak tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal hapus kontak")
	}

	return utils.SuccessResponse(c, "kontak berhasil dihapus", nil)
}
