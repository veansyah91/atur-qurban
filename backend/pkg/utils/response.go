package utils

import "github.com/gofiber/fiber/v2"

// BaseResponse format respons standar untuk semua endpoint
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Meta informasi pagination
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse format respons untuk data dengan pagination
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

// SuccessResponse mengirim respons sukses dengan status 200
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirim respons error dengan kode status yang ditentukan
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(BaseResponse{
		Success: false,
		Message: message,
	})
}

// PaginateResponse mengirim respons sukses dengan data dan informasi pagination
func PaginateResponse(c *fiber.Ctx, message string, data interface{}, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
