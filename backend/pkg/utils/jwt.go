package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims struktur claim untuk JWT
type CustomClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Type    string `json:"type"` // "access" atau "refresh"
	JTI     string `json:"jti"`  // JWT ID untuk blacklist
	jwt.RegisteredClaims
}

// GenerateTokenPair membuat access token dan refresh token
func GenerateTokenPair(userID string, isAdmin bool, jwtSecret string, accessExpiryHour, refreshExpiryDay int) (accessToken, refreshToken string, err error) {
	now := time.Now()
	jti := uuid.New().String()

	// Access Token
	accessClaims := &CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		Type:    "access",
		JTI:     jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accessExpiryHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("GenerateTokenPair: gagal sign access token: %w", err)
	}

	// Refresh Token
	refreshClaims := &CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		Type:    "refresh",
		JTI:     jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(refreshExpiryDay) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("GenerateTokenPair: gagal sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ParseToken memparse dan validasi token JWT
func ParseToken(tokenStr, jwtSecret string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ParseToken: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("ParseToken: token tidak valid")
	}

	return claims, nil
}

// GetRemainingTime menghitung sisa waktu token dalam detik
func GetRemainingTime(claims *CustomClaims) int64 {
	if claims.ExpiresAt == nil {
		return 0
	}
	remaining := claims.ExpiresAt.Time.Unix() - time.Now().Unix()
	if remaining < 0 {
		return 0
	}
	return remaining
}
