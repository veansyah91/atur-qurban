package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/pkg/utils"
)

// AuthHandler menangani request authentication
type AuthHandler struct {
	authService service.AuthService
	cfg         *config.AppConfig
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(authService service.AuthService, cfg *config.AppConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// request structs

type RegisterRequest struct {
Name            string `json:"name"`
Phone           string `json:"phone"`
Password        string `json:"password"`
ConfirmPassword string `json:"confirm_password"`
}

type VerifyOTPRequest struct {
Phone string `json:"phone"`
OTP   string `json:"otp"`
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

type ResendOTPRequest struct {
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

// Helper untuk set cookie
func (h *AuthHandler) setTokenCookies(c *fiber.Ctx, accessToken, refreshToken string) {
	cookieSecure := h.cfg.App.Env == "production"

	// Access Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Duration(h.cfg.JWT.ExpiryHour) * time.Hour),
		HTTPOnly: true,
		Secure:   cookieSecure,
		SameSite: "Lax",
	})

	// Refresh Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(h.cfg.JWT.RefreshExpiryDay) * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   cookieSecure,
		SameSite: "Lax",
	})
}

// Helper untuk hapus cookie
func (h *AuthHandler) clearTokenCookies(c *fiber.Ctx) {
	cookieSecure := h.cfg.App.Env == "production"

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   cookieSecure,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   cookieSecure,
		SameSite: "Lax",
	})
}

// @Summary     Resend OTP registrasi
// @Description Mengirim ulang OTP untuk proses registrasi yang sedang pending
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body ForgotPasswordRequest true "Nomor telepon"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/register/resend-otp [post]
// ResendRegisterOTP handler — POST /api/v1/auth/register/resend-otp
func (h *AuthHandler) ResendRegisterOTP(c *fiber.Ctx) error {
	var req ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Phone == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "phone harus diisi")
	}

	otp, err := h.authService.ResendRegisterOTP(c.Context(), req.Phone)
	if err != nil {
		if errors.Is(err, service.ErrNoPendingRegistration) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "tidak ada registrasi pending untuk nomor ini, silakan daftar ulang")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal resend OTP")
	}

	resp := fiber.Map{}
	if h.cfg.App.Env != "production" && otp != "" {
		resp["otp"] = otp
	}
	return utils.SuccessResponse(c, "OTP berhasil dikirim ulang", resp)
}

// @Summary     Request OTP untuk registrasi
// @Description Mengirim OTP ke nomor telepon untuk proses registrasi pengguna
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body RegisterRequest true "Data registrasi dengan name, phone, password, dan confirm_password"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     409 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/register/request-otp [post]
// RequestRegisterOTP handler — POST /api/v1/auth/register/request-otp
func (h *AuthHandler) RequestRegisterOTP(c *fiber.Ctx) error {
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
	if req.ConfirmPassword == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "confirm_password harus diisi")
	}
	if req.Password != req.ConfirmPassword {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "password dan confirm_password tidak cocok")
	}

	otp, err := h.authService.RequestRegisterOTP(c.Context(), req.Name, req.Phone, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrPhoneAlreadyRegistered) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "nomor telepon sudah terdaftar")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal request OTP: "+err.Error())
	}

	resp := fiber.Map{}
	if h.cfg.App.Env != "production" && otp != "" {
		resp["otp"] = otp
	}
	return utils.SuccessResponse(c, "OTP berhasil dikirim", resp)
}

// @Summary     Verifikasi OTP registrasi
// @Description Memverifikasi OTP dan menyelesaikan proses registrasi pengguna
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body VerifyOTPRequest true "Nomor telepon dan OTP"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/register/verify [post]
// VerifyRegisterOTP handler — POST /api/v1/auth/register/verify
func (h *AuthHandler) VerifyRegisterOTP(c *fiber.Ctx) error {
	var req VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Phone == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "phone harus diisi")
	}
	if req.OTP == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "otp harus diisi")
	}

	userResp, accessToken, refreshToken, err := h.authService.VerifyRegisterOTP(c.Context(), req.Phone, req.OTP)
	if err != nil {
		if errors.Is(err, service.ErrOTPExpiredCanResend) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP sudah expired, gunakan /register/resend-otp untuk mendapatkan OTP baru")
		}
		if errors.Is(err, service.ErrOTPExpired) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP tidak valid atau sudah expired")
		}
		if errors.Is(err, service.ErrOTPInvalid) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP salah")
		}
		if errors.Is(err, service.ErrPendingDataExpired) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "data registrasi sudah expired, silakan daftar ulang")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal verifikasi OTP")
	}

	h.setTokenCookies(c, accessToken, refreshToken)

	return utils.SuccessResponse(c, "registrasi berhasil", fiber.Map{
		"user": userResp,
	})
}

// @Summary     Endpoint deprecated
// @Description Endpoint registrasi lama telah deprecated. Gunakan /register/request-otp kemudian /register/verify
// @Tags        auth
// @Accept      json
// @Produce     json
// @Success     410 {object} utils.BaseResponse
// @Router      /api/v1/auth/register [post]
// Register handler — POST /api/v1/auth/register (deprecated, use request-otp then verify)
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusGone, "endpoint ini deprecated, gunakan /register/request-otp kemudian /register/verify")
}

// @Summary     Login pengguna admin
// @Description Login dengan nomor telepon dan password untuk pengguna admin
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body LoginRequest true "Nomor telepon dan password"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/login [post]
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
		if errors.Is(err, service.ErrUserNotFound) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "phone atau password salah")
		}
		if errors.Is(err, service.ErrPasswordIncorrect) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "phone atau password salah")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal login: "+err.Error())
	}

	h.setTokenCookies(c, accessToken, refreshToken)

	return utils.SuccessResponse(c, "login berhasil", fiber.Map{
		"user": userResp,
	})
}

// @Summary     Refresh token JWT
// @Description Membuat token JWT baru menggunakan refresh token yang valid
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body RefreshTokenRequest true "Refresh token"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/refresh [post]
// RefreshToken handler — POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "refresh_token tidak ditemukan di cookie")
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrTokenBlacklisted) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "token sudah invalid")
		}
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "gagal refresh token: "+err.Error())
	}

	h.setTokenCookies(c, newAccessToken, newRefreshToken)

	return utils.SuccessResponse(c, "refresh token berhasil", nil)
}

// @Summary     Logout pengguna
// @Description Logout dan invalidkan access token pengguna
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/logout [post]
// Logout handler — POST /api/v1/auth/logout (protected)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		// Fallback to Header for flexibility during transition
		token = c.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token != "" {
		if err := h.authService.Logout(c.Context(), token); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal logout: "+err.Error())
		}
	}

	h.clearTokenCookies(c)

	return utils.SuccessResponse(c, "logout berhasil", nil)
}

// @Summary     Ambil data user saat ini
// @Description Mengambil data user yang sedang login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/me [get]
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

// @Summary     Forgot password request
// @Description Mengirim OTP ke nomor telepon untuk reset password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body ForgotPasswordRequest true "Nomor telepon"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/forgot-password [post]
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
		if errors.Is(err, service.ErrUserNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "nomor telepon tidak terdaftar")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal proses forgot password: "+err.Error())
	}

	resp := fiber.Map{}
	if h.cfg.App.Env != "production" && otp != "" {
		resp["otp"] = otp
	}
	return utils.SuccessResponse(c, "OTP berhasil dikirim", resp)
}

// @Summary     Reset password
// @Description Mereset password pengguna menggunakan OTP yang valid
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body ResetPasswordRequest true "Nomor telepon, OTP, dan password baru"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     404 {object} utils.BaseResponse
// @Failure     500 {object} utils.BaseResponse
// @Router      /api/v1/auth/reset-password [post]
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
		if errors.Is(err, service.ErrOTPExpired) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP tidak valid atau sudah expired")
		}
		if errors.Is(err, service.ErrOTPInvalid) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "OTP salah")
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user tidak ditemukan")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gagal reset password: "+err.Error())
	}

	return utils.SuccessResponse(c, "password berhasil direset", nil)
}
