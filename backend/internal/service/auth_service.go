package service

import (
"context"
"encoding/json"
"errors"
"fmt"
"io"
"log"
"math/rand"
"time"

"github.com/redis/go-redis/v9"
"github.com/username/qurban-app/config"
"github.com/username/qurban-app/internal/model"
"github.com/username/qurban-app/internal/repository"
"github.com/username/qurban-app/pkg/utils"
"golang.org/x/crypto/bcrypt"
)

// Sentinel errors untuk AuthService
var (
	ErrPhoneAlreadyRegistered  = errors.New("nomor telepon sudah terdaftar")
	ErrOTPInvalid              = errors.New("OTP tidak valid")
	ErrOTPExpired              = errors.New("OTP tidak ditemukan atau sudah expired")
	ErrOTPExpiredCanResend     = errors.New("OTP sudah expired, silakan kirim ulang OTP")
	ErrPendingDataExpired      = errors.New("data registrasi tidak ditemukan atau sudah expired")
	ErrNoPendingRegistration   = errors.New("tidak ada data registrasi pending untuk nomor ini")
	ErrUserNotFound            = errors.New("user tidak ditemukan")
	ErrPasswordMissing         = errors.New("user belum memiliki password")
	ErrPasswordIncorrect       = errors.New("password salah")
	ErrTokenBlacklisted        = errors.New("token sudah di-blacklist")
)

// AuthService interface untuk operasi authentication
type AuthService interface {
	RequestRegisterOTP(ctx context.Context, name, phone, password string) (string, error)
	ResendRegisterOTP(ctx context.Context, phone string) (string, error)
	VerifyRegisterOTP(ctx context.Context, phone, otp string) (*model.UserResponse, string, string, error)
	Login(phone, password string) (*model.UserResponse, string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	ValidateToken(token string) (*utils.CustomClaims, error)
	Logout(ctx context.Context, token string) error
	GetUserByID(userID string) (*model.UserResponse, error)
	ForgotPassword(ctx context.Context, phone string) (string, error)
	ResetPassword(ctx context.Context, phone, otp, newPassword string) error
}

// authService implementasi AuthService
type authService struct {
	userRepo            repository.UserRepository
	rdb                 *redis.Client
	cfg                 *config.AppConfig
	notificationService NotificationService
	logger              *log.Logger
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(userRepo repository.UserRepository, rdb *redis.Client, cfg *config.AppConfig, notifService NotificationService, logger *log.Logger) AuthService {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &authService{
		userRepo:            userRepo,
		rdb:                 rdb,
		cfg:                 cfg,
		notificationService: notifService,
		logger:              logger,
	}
}

// RequestRegisterOTP menyimpan data registrasi pending ke Redis dan mengirim OTP
func (s *authService) RequestRegisterOTP(ctx context.Context, name, phone, password string) (string, error) {
	phone = utils.NormalizePhoneNumber(phone)

	// Cek apakah nomor sudah terdaftar
	existingUser, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", fmt.Errorf("RequestRegisterOTP: %w", err)
	}
	if existingUser != nil {
		return "", fmt.Errorf("RequestRegisterOTP: %w", ErrPhoneAlreadyRegistered)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("RequestRegisterOTP: gagal hash password: %w", err)
	}

	// Simpan data pending ke Redis (TTL 10 menit)
	pendingKey := fmt.Sprintf("register:pending:%s", phone)
	pendingData := fmt.Sprintf(`{"name":"%s","password":"%s"}`, name, string(hashedPassword))
	if err := s.rdb.Set(ctx, pendingKey, pendingData, 10*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("RequestRegisterOTP: gagal simpan data pending: %w", err)
	}

	// Generate OTP 6 digit
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Hash OTP dengan bcrypt sebelum disimpan
	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("RequestRegisterOTP: gagal hash OTP: %w", err)
	}

	// Simpan OTP hash ke Redis (TTL 5 menit)
	otpKey := fmt.Sprintf("otp:register:%s", phone)
	if err := s.rdb.Set(ctx, otpKey, string(hashedOTP), 5*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("RequestRegisterOTP: gagal simpan OTP: %w", err)
	}

	// Kirim OTP via WhatsApp
	if s.notificationService != nil {
		message := fmt.Sprintf("Kode OTP registrasi Anda: %s. Jangan bagikan kode ini ke siapapun. Berlaku 5 menit.", otp)
		if err := s.notificationService.Send(ctx, phone, message); err != nil {
			s.logger.Printf("RequestRegisterOTP: gagal mengirim OTP via WhatsApp: %v\n", err)
		}
	}

	return otp, nil
}

// ResendRegisterOTP mengirim ulang OTP untuk registrasi yang pending
func (s *authService) ResendRegisterOTP(ctx context.Context, phone string) (string, error) {
	phone = utils.NormalizePhoneNumber(phone)

	// Cek apakah data pending masih ada di Redis
	pendingKey := fmt.Sprintf("register:pending:%s", phone)
	if _, err := s.rdb.Get(ctx, pendingKey).Result(); err != nil {
		return "", fmt.Errorf("ResendRegisterOTP: %w", ErrNoPendingRegistration)
	}

	// Generate OTP baru
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Hash OTP dengan bcrypt sebelum disimpan
	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ResendRegisterOTP: gagal hash OTP: %w", err)
	}

	// Simpan OTP baru ke Redis (reset TTL 5 menit)
	otpKey := fmt.Sprintf("otp:register:%s", phone)
	if err := s.rdb.Set(ctx, otpKey, string(hashedOTP), 5*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("ResendRegisterOTP: gagal simpan OTP: %w", err)
	}

	// Perpanjang TTL data pending menjadi 10 menit
	if err := s.rdb.Expire(ctx, pendingKey, 10*time.Minute).Err(); err != nil {
		s.logger.Printf("ResendRegisterOTP: gagal perpanjang TTL data pending: %v\n", err)
	}

	// Kirim OTP baru via WhatsApp
	if s.notificationService != nil {
		message := fmt.Sprintf("Kode OTP registrasi Anda: %s. Jangan bagikan kode ini ke siapapun. Berlaku 5 menit.", otp)
		if err := s.notificationService.Send(ctx, phone, message); err != nil {
			s.logger.Printf("ResendRegisterOTP: gagal mengirim OTP via WhatsApp: %v\n", err)
		}
	}

	return otp, nil
}

// VerifyRegisterOTP memvalidasi OTP dan membuat user baru dengan verified_at diisi
func (s *authService) VerifyRegisterOTP(ctx context.Context, phone, otp string) (*model.UserResponse, string, string, error) {
	phone = utils.NormalizePhoneNumber(phone)

	// Ambil OTP hash dari Redis
	otpKey := fmt.Sprintf("otp:register:%s", phone)
	storedOTPHash, err := s.rdb.Get(ctx, otpKey).Result()
	if err != nil {
		// OTP expired — cek apakah data pending masih ada
		pendingKey := fmt.Sprintf("register:pending:%s", phone)
		if _, pendingErr := s.rdb.Get(ctx, pendingKey).Result(); pendingErr == nil {
			return nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", ErrOTPExpiredCanResend)
		}
		return nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", ErrPendingDataExpired)
	}

// Validasi OTP dengan bcrypt
if err := bcrypt.CompareHashAndPassword([]byte(storedOTPHash), []byte(otp)); err != nil {
return nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", ErrOTPInvalid)
}

// Ambil data pending dari Redis
pendingKey := fmt.Sprintf("register:pending:%s", phone)
pendingData, err := s.rdb.Get(ctx, pendingKey).Result()
if err != nil {
return nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", ErrPendingDataExpired)
}

// Parse data pending (JSON: {name, password})
var pending map[string]string
if err := json.Unmarshal([]byte(pendingData), &pending); err != nil {
return nil, "", "", fmt.Errorf("VerifyRegisterOTP: gagal parse data pending: %w", err)
}

name := pending["name"]
hashedPassword := pending["password"]

// Buat user baru dengan verified_at diisi
now := time.Now()
user := &model.User{
Name:       name,
Phone:      phone,
IsAdmin:    false,
Password:   &hashedPassword,
VerifiedAt: &now,
}

if err := s.userRepo.Create(user); err != nil {
return nil, "", "", fmt.Errorf("VerifyRegisterOTP: %w", err)
}

// Generate token pair
accessToken, refreshToken, err := utils.GenerateTokenPair(
user.ID, user.IsAdmin,
s.cfg.JWT.Secret, s.cfg.JWT.ExpiryHour, s.cfg.JWT.RefreshExpiryDay,
)
if err != nil {
return nil, "", "", fmt.Errorf("VerifyRegisterOTP: GenerateTokenPair error: %w", err)
}

// Hapus OTP dan data pending dari Redis
s.rdb.Del(ctx, otpKey, pendingKey)

return user.ToResponse(), accessToken, refreshToken, nil
}

// Login melakukan login user berdasarkan phone + password
func (s *authService) Login(phone, password string) (*model.UserResponse, string, string, error) {
phone = utils.NormalizePhoneNumber(phone)

user, err := s.userRepo.FindByPhone(phone)
if err != nil {
return nil, "", "", fmt.Errorf("Login: %w", err)
}
if user == nil {
return nil, "", "", fmt.Errorf("Login: %w", ErrUserNotFound)
}

// Verifikasi password
if user.Password == nil {
return nil, "", "", fmt.Errorf("Login: %w", ErrPasswordMissing)
}
if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
return nil, "", "", fmt.Errorf("Login: %w", ErrPasswordIncorrect)
}

if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
return nil, "", "", fmt.Errorf("Login: UpdateLastLogin error: %w", err)
}

accessToken, refreshToken, err := utils.GenerateTokenPair(
user.ID, user.IsAdmin,
s.cfg.JWT.Secret, s.cfg.JWT.ExpiryHour, s.cfg.JWT.RefreshExpiryDay,
)
if err != nil {
return nil, "", "", fmt.Errorf("Login: GenerateTokenPair error: %w", err)
}

return user.ToResponse(), accessToken, refreshToken, nil
}

// RefreshToken membuat token pair baru dari refresh token yang valid
func (s *authService) RefreshToken(refreshTokenStr string) (string, string, error) {
claims, err := utils.ParseToken(refreshTokenStr, s.cfg.JWT.Secret)
if err != nil {
return "", "", fmt.Errorf("RefreshToken: invalid token: %w", err)
}
if claims.Type != "refresh" {
return "", "", fmt.Errorf("RefreshToken: token bukan refresh token")
}

ctx := context.Background()
blacklistKey := fmt.Sprintf("blacklist:%s", claims.JTI)
exists, err := s.rdb.Exists(ctx, blacklistKey).Result()
if err != nil {
return "", "", fmt.Errorf("RefreshToken: redis error: %w", err)
}
if exists > 0 {
return "", "", fmt.Errorf("RefreshToken: %w", ErrTokenBlacklisted)
}

user, err := s.userRepo.FindByID(claims.UserID)
if err != nil {
return "", "", fmt.Errorf("RefreshToken: %w", err)
}
if user == nil {
return "", "", fmt.Errorf("RefreshToken: user tidak ditemukan")
}

// Blacklist token lama
ttl := utils.GetRemainingTime(claims)
if ttl > 0 {
if err := s.rdb.Set(ctx, blacklistKey, "1", time.Duration(ttl)*time.Second).Err(); err != nil {
return "", "", fmt.Errorf("RefreshToken: gagal blacklist token lama: %w", err)
}
}

return utils.GenerateTokenPair(
user.ID, user.IsAdmin,
s.cfg.JWT.Secret, s.cfg.JWT.ExpiryHour, s.cfg.JWT.RefreshExpiryDay,
)
}

// ValidateToken memvalidasi access token
func (s *authService) ValidateToken(token string) (*utils.CustomClaims, error) {
claims, err := utils.ParseToken(token, s.cfg.JWT.Secret)
if err != nil {
return nil, fmt.Errorf("ValidateToken: %w", err)
}
if claims.Type != "access" {
return nil, fmt.Errorf("ValidateToken: token bukan access token")
}

ctx := context.Background()
blacklistKey := fmt.Sprintf("blacklist:%s", claims.JTI)
exists, err := s.rdb.Exists(ctx, blacklistKey).Result()
if err != nil {
return nil, fmt.Errorf("ValidateToken: redis error: %w", err)
}
if exists > 0 {
return nil, fmt.Errorf("ValidateToken: %w", ErrTokenBlacklisted)
}

return claims, nil
}

// Logout menambahkan token ke blacklist
func (s *authService) Logout(ctx context.Context, token string) error {
claims, err := utils.ParseToken(token, s.cfg.JWT.Secret)
if err != nil {
return fmt.Errorf("Logout: invalid token: %w", err)
}

blacklistKey := fmt.Sprintf("blacklist:%s", claims.JTI)
ttl := utils.GetRemainingTime(claims)
if ttl > 0 {
if err := s.rdb.Set(ctx, blacklistKey, "1", time.Duration(ttl)*time.Second).Err(); err != nil {
return fmt.Errorf("Logout: gagal blacklist token: %w", err)
}
}

return nil
}

// GetUserByID mengambil data user berdasarkan ID
func (s *authService) GetUserByID(userID string) (*model.UserResponse, error) {
user, err := s.userRepo.FindByID(userID)
if err != nil {
return nil, fmt.Errorf("GetUserByID: %w", err)
}
if user == nil {
return nil, fmt.Errorf("GetUserByID: user tidak ditemukan")
}
return user.ToResponse(), nil
}

// ForgotPassword membuat OTP reset password dan menyimpannya di Redis
func (s *authService) ForgotPassword(ctx context.Context, phone string) (string, error) {
phone = utils.NormalizePhoneNumber(phone)

user, err := s.userRepo.FindByPhone(phone)
if err != nil {
return "", fmt.Errorf("ForgotPassword: %w", err)
}
if user == nil {
return "", fmt.Errorf("ForgotPassword: %w", ErrUserNotFound)
}

// Generate OTP 6 digit
otp := fmt.Sprintf("%06d", rand.Intn(1000000))

// Hash OTP dengan bcrypt sebelum disimpan
hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
if err != nil {
return "", fmt.Errorf("ForgotPassword: gagal hash OTP: %w", err)
}

// Simpan ke Redis dengan TTL 5 menit
otpKey := fmt.Sprintf("otp:reset:%s", phone)
if err := s.rdb.Set(ctx, otpKey, string(hashedOTP), 5*time.Minute).Err(); err != nil {
return "", fmt.Errorf("ForgotPassword: gagal simpan OTP: %w", err)
}

// Kirim OTP via WhatsApp menggunakan notification service
if s.notificationService != nil {
message := fmt.Sprintf("Kode OTP Anda: %s. Jangan bagikan kode ini ke siapapun. Berlaku 5 menit.", otp)
if err := s.notificationService.Send(ctx, phone, message); err != nil {
// Log error tapi jangan return - OTP tetap disimpan di Redis
s.logger.Printf("ForgotPassword: gagal mengirim OTP via WhatsApp: %v\n", err)
}
}

// Kembalikan OTP hanya jika environment bukan production
if s.cfg.App.Env == "production" {
return "", nil
}
return otp, nil
}

// ResetPassword memvalidasi OTP dan mengupdate password baru
func (s *authService) ResetPassword(ctx context.Context, phone, otp, newPassword string) error {
phone = utils.NormalizePhoneNumber(phone)

// Ambil OTP hash dari Redis
otpKey := fmt.Sprintf("otp:reset:%s", phone)
storedOTPHash, err := s.rdb.Get(ctx, otpKey).Result()
if err != nil {
return fmt.Errorf("ResetPassword: %w", ErrOTPExpired)
}

// Validasi OTP dengan bcrypt
if err := bcrypt.CompareHashAndPassword([]byte(storedOTPHash), []byte(otp)); err != nil {
return fmt.Errorf("ResetPassword: %w", ErrOTPInvalid)
}

// Cari user
user, err := s.userRepo.FindByPhone(phone)
if err != nil {
return fmt.Errorf("ResetPassword: %w", err)
}
if user == nil {
return fmt.Errorf("ResetPassword: %w", ErrUserNotFound)
}

// Hash password baru
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
if err != nil {
return fmt.Errorf("ResetPassword: gagal hash password: %w", err)
}

// Update password di DB
if err := s.userRepo.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
return fmt.Errorf("ResetPassword: %w", err)
}

// Hapus OTP dari Redis
s.rdb.Del(ctx, otpKey)

return nil
}
