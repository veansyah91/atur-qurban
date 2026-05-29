package service

import (
"context"
"fmt"
"math/rand"
"time"

"github.com/redis/go-redis/v9"
"github.com/username/qurban-app/config"
"github.com/username/qurban-app/internal/model"
"github.com/username/qurban-app/internal/repository"
"github.com/username/qurban-app/pkg/utils"
"golang.org/x/crypto/bcrypt"
)

// AuthService interface untuk operasi authentication
type AuthService interface {
Register(name, phone, password string) (*model.UserResponse, string, string, error)
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
userRepo repository.UserRepository
rdb      *redis.Client
cfg      *config.AppConfig
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(userRepo repository.UserRepository, rdb *redis.Client, cfg *config.AppConfig) AuthService {
return &authService{
userRepo: userRepo,
rdb:      rdb,
cfg:      cfg,
}
}

// Register membuat user baru dengan password
func (s *authService) Register(name, phone, password string) (*model.UserResponse, string, string, error) {
phone = utils.NormalizePhoneNumber(phone)

existingUser, err := s.userRepo.FindByPhone(phone)
if err != nil {
return nil, "", "", fmt.Errorf("Register: FindByPhone error: %w", err)
}
if existingUser != nil {
return nil, "", "", fmt.Errorf("Register: nomor telepon sudah terdaftar")
}

// Hash password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
return nil, "", "", fmt.Errorf("Register: gagal hash password: %w", err)
}

hashed := string(hashedPassword)
user := &model.User{
Name:     name,
Phone:    phone,
IsAdmin:  false,
Password: &hashed,
}

if err := s.userRepo.Create(user); err != nil {
return nil, "", "", fmt.Errorf("Register: %w", err)
}

accessToken, refreshToken, err := utils.GenerateTokenPair(
user.ID, user.IsAdmin,
s.cfg.JWT.Secret, s.cfg.JWT.ExpiryHour, s.cfg.JWT.RefreshExpiryDay,
)
if err != nil {
return nil, "", "", fmt.Errorf("Register: GenerateTokenPair error: %w", err)
}

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
return nil, "", "", fmt.Errorf("Login: user tidak ditemukan")
}

// Verifikasi password
if user.Password == nil {
return nil, "", "", fmt.Errorf("Login: user belum memiliki password")
}
if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
return nil, "", "", fmt.Errorf("Login: password salah")
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
return "", "", fmt.Errorf("RefreshToken: token sudah di-blacklist")
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
return nil, fmt.Errorf("ValidateToken: token sudah di-blacklist")
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
return "", fmt.Errorf("ForgotPassword: user tidak ditemukan")
}

// Generate OTP 6 digit
otp := fmt.Sprintf("%06d", rand.Intn(1000000))

// Simpan ke Redis dengan TTL 5 menit
otpKey := fmt.Sprintf("otp:reset:%s", phone)
if err := s.rdb.Set(ctx, otpKey, otp, 5*time.Minute).Err(); err != nil {
return "", fmt.Errorf("ForgotPassword: gagal simpan OTP: %w", err)
}

// Di production, kirim OTP via WhatsApp/SMS
// Untuk development, OTP dikembalikan langsung dalam response
return otp, nil
}

// ResetPassword memvalidasi OTP dan mengupdate password baru
func (s *authService) ResetPassword(ctx context.Context, phone, otp, newPassword string) error {
phone = utils.NormalizePhoneNumber(phone)

// Ambil OTP dari Redis
otpKey := fmt.Sprintf("otp:reset:%s", phone)
storedOTP, err := s.rdb.Get(ctx, otpKey).Result()
if err != nil {
return fmt.Errorf("ResetPassword: OTP tidak ditemukan atau sudah expired")
}

// Bandingkan OTP
if storedOTP != otp {
return fmt.Errorf("ResetPassword: OTP tidak valid")
}

// Cari user
user, err := s.userRepo.FindByPhone(phone)
if err != nil {
return fmt.Errorf("ResetPassword: %w", err)
}
if user == nil {
return fmt.Errorf("ResetPassword: user tidak ditemukan")
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
