package handler

import (
"github.com/gofiber/fiber/v2"
"github.com/username/qurban-app/internal/model"
"github.com/username/qurban-app/internal/service"
"github.com/username/qurban-app/pkg/utils"
)

// AuthHandler menangani request authentication
type AuthHandler struct {
authService service.AuthService
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
return &AuthHandler{authService: authService}
}

// request structs

type RegisterRequest struct {
Name     string `json:"name"`
Phone    string `json:"phone"`
Password string `json:"password"`
}

type LoginRequest struct {
Phone    string `json:"phone"`
Password string `json:"password"`
}

type RefreshTokenRequest struct {
RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
Phone string `json:"phone"`
}

type ResetPasswordRequest struct {
Phone       string `json:"phone"`
OTP         string `json:"otp"`
NewPassword string `json:"new_password"`
}

// TokenResponse untuk response token
type TokenResponse struct {
User         *model.UserResponse `json:"user"`
AccessToken  string              `json:"access_token"`
RefreshToken string              `json:"refresh_token"`
}

// Register handler — POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
var req RegisterRequest
if err := c.BodyParser(&req); err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
}

if req.Name == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "name harus diisi")
}
if req.Phone == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "phone harus diisi")
}
if req.Password == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "password harus diisi")
}

userResp, accessToken, refreshToken, err := h.authService.Register(req.Name, req.Phone, req.Password)
if err != nil {
if err.Error() == "Register: nomor telepon sudah terdaftar" {
return utils.ErrorResponse(c, fiber.StatusConflict, "nomor telepon sudah terdaftar")
}
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal register: "+err.Error())
}

return utils.SuccessResponse(c, "register berhasil", fiber.Map{
"user":          userResp,
"access_token":  accessToken,
"refresh_token": refreshToken,
})
}

// Login handler — POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
var req LoginRequest
if err := c.BodyParser(&req); err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
}

if req.Phone == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "phone harus diisi")
}
if req.Password == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "password harus diisi")
}

userResp, accessToken, refreshToken, err := h.authService.Login(req.Phone, req.Password)
if err != nil {
switch err.Error() {
case "Login: user tidak ditemukan":
return utils.ErrorResponse(c, fiber.StatusUnauthorized, "phone atau password salah")
case "Login: password salah":
return utils.ErrorResponse(c, fiber.StatusUnauthorized, "phone atau password salah")
default:
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal login: "+err.Error())
}
}

return utils.SuccessResponse(c, "login berhasil", fiber.Map{
"user":          userResp,
"access_token":  accessToken,
"refresh_token": refreshToken,
})
}

// RefreshToken handler — POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
var req RefreshTokenRequest
if err := c.BodyParser(&req); err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
}

if req.RefreshToken == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "refresh_token harus diisi")
}

newAccessToken, newRefreshToken, err := h.authService.RefreshToken(req.RefreshToken)
if err != nil {
if err.Error() == "RefreshToken: token sudah di-blacklist" {
return utils.ErrorResponse(c, fiber.StatusUnauthorized, "token sudah invalid")
}
return utils.ErrorResponse(c, fiber.StatusUnauthorized, "gagal refresh token: "+err.Error())
}

return utils.SuccessResponse(c, "refresh token berhasil", fiber.Map{
"access_token":  newAccessToken,
"refresh_token": newRefreshToken,
})
}

// Logout handler — POST /api/v1/auth/logout (protected)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
token := c.Get("Authorization")
if token == "" {
return utils.ErrorResponse(c, fiber.StatusUnauthorized, "authorization header tidak ditemukan")
}
if len(token) > 7 && token[:7] == "Bearer " {
token = token[7:]
}

if err := h.authService.Logout(c.Context(), token); err != nil {
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal logout: "+err.Error())
}

return utils.SuccessResponse(c, "logout berhasil", nil)
}

// Me handler — GET /api/v1/auth/me (protected)
func (h *AuthHandler) Me(c *fiber.Ctx) error {
userID := c.Locals("user_id")
if userID == nil {
return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_id tidak ditemukan")
}

userResp, err := h.authService.GetUserByID(userID.(string))
if err != nil {
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal ambil data user: "+err.Error())
}

return utils.SuccessResponse(c, "data user berhasil diambil", userResp)
}

// ForgotPassword handler — POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
var req ForgotPasswordRequest
if err := c.BodyParser(&req); err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
}

if req.Phone == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "phone harus diisi")
}

otp, err := h.authService.ForgotPassword(c.Context(), req.Phone)
if err != nil {
if err.Error() == "ForgotPassword: user tidak ditemukan" {
return utils.ErrorResponse(c, fiber.StatusNotFound, "nomor telepon tidak terdaftar")
}
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal proses forgot password: "+err.Error())
}

// Di production, OTP dikirim via SMS/WhatsApp — tidak dikembalikan dalam response
// Untuk development, OTP dikembalikan agar mudah testing
return utils.SuccessResponse(c, "OTP berhasil dikirim", fiber.Map{
"otp": otp, // hapus field ini di production
})
}

// ResetPassword handler — POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
var req ResetPasswordRequest
if err := c.BodyParser(&req); err != nil {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
}

if req.Phone == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "phone harus diisi")
}
if req.OTP == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "otp harus diisi")
}
if req.NewPassword == "" {
return utils.ErrorResponse(c, fiber.StatusBadRequest, "new_password harus diisi")
}

if err := h.authService.ResetPassword(c.Context(), req.Phone, req.OTP, req.NewPassword); err != nil {
switch err.Error() {
case "ResetPassword: OTP tidak ditemukan atau sudah expired":
return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP tidak valid atau sudah expired")
case "ResetPassword: OTP tidak valid":
return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP salah")
case "ResetPassword: user tidak ditemukan":
return utils.ErrorResponse(c, fiber.StatusNotFound, "user tidak ditemukan")
default:
return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal reset password: "+err.Error())
}
}

return utils.SuccessResponse(c, "password berhasil direset", nil)
}
