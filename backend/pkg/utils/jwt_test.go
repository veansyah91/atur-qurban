package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateTokenPair(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - menghasilkan access + refresh token yang valid", func(t *testing.T) {
		t.Parallel()

		userID := "user-123"
		isAdmin := false
		jwtSecret := "test-secret-key-very-secure"
		accessExpiryHour := 1
		refreshExpiryDay := 7

		accessToken, refreshToken, err := GenerateTokenPair(
			userID, isAdmin, jwtSecret, accessExpiryHour, refreshExpiryDay,
		)

		require.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		// Parse access token
		accessClaims, err := ParseToken(accessToken, jwtSecret)
		require.NoError(t, err)
		assert.Equal(t, userID, accessClaims.UserID)
		assert.Equal(t, isAdmin, accessClaims.IsAdmin)
		assert.Equal(t, "access", accessClaims.Type)
		assert.NotEmpty(t, accessClaims.JTI)

		// Parse refresh token
		refreshClaims, err := ParseToken(refreshToken, jwtSecret)
		require.NoError(t, err)
		assert.Equal(t, userID, refreshClaims.UserID)
		assert.Equal(t, isAdmin, refreshClaims.IsAdmin)
		assert.Equal(t, "refresh", refreshClaims.Type)
		assert.Equal(t, accessClaims.JTI, refreshClaims.JTI)
	})

	t.Run("Sukses - admin flag diteruskan ke token", func(t *testing.T) {
		t.Parallel()

		accessToken, _, err := GenerateTokenPair(
			"admin-user", true, "secret", 1, 7,
		)

		require.NoError(t, err)
		claims, err := ParseToken(accessToken, "secret")
		require.NoError(t, err)
		assert.True(t, claims.IsAdmin)
	})
}

func TestParseToken(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - parse valid token", func(t *testing.T) {
		t.Parallel()

		accessToken, _, err := GenerateTokenPair(
			"user-123", false, "secret", 1, 7,
		)
		require.NoError(t, err)

		claims, err := ParseToken(accessToken, "secret")
		require.NoError(t, err)

		assert.Equal(t, "user-123", claims.UserID)
		assert.False(t, claims.IsAdmin)
		assert.Equal(t, "access", claims.Type)
	})

	t.Run("Gagal - token dengan signature salah", func(t *testing.T) {
		t.Parallel()

		accessToken, _, err := GenerateTokenPair(
			"user-123", false, "secret-1", 1, 7,
		)
		require.NoError(t, err)

		_, err = ParseToken(accessToken, "secret-2") // Secret berbeda
		assert.Error(t, err)
	})

	t.Run("Gagal - token malformed", func(t *testing.T) {
		t.Parallel()

		_, err := ParseToken("invalid.token.format", "secret")
		assert.Error(t, err)
	})

	t.Run("Gagal - token kosong", func(t *testing.T) {
		t.Parallel()

		_, err := ParseToken("", "secret")
		assert.Error(t, err)
	})

	t.Run("Gagal - token expired", func(t *testing.T) {
		t.Parallel()

		// Buat token dengan expiry 0 detik
		userID := "user-123"
		isAdmin := false
		jwtSecret := "secret"

		// Gunakan GenerateTokenPair dengan expiry hour -1 (sudah expired)
		accessToken, _, err := GenerateTokenPair(
			userID, isAdmin, jwtSecret, -1, 7,
		)
		require.NoError(t, err)

		// Berikan delay kecil agar token benar-benar expired
		time.Sleep(100 * time.Millisecond)

		_, err = ParseToken(accessToken, jwtSecret)
		assert.Error(t, err)
	})
}

func TestGetRemainingTime(t *testing.T) {
	t.Parallel()

	t.Run("Sukses - token dengan sisa waktu > 0", func(t *testing.T) {
		t.Parallel()

		accessToken, _, err := GenerateTokenPair(
			"user-123", false, "secret", 24, 7, // 24 jam expiry
		)
		require.NoError(t, err)

		claims, err := ParseToken(accessToken, "secret")
		require.NoError(t, err)

		remaining := GetRemainingTime(claims)
		assert.Greater(t, remaining, int64(0))
		assert.LessOrEqual(t, remaining, int64(24*3600)) // Kurang dari atau sama 24 jam
	})

	t.Run("Gagal - token sudah expired mengembalikan 0", func(t *testing.T) {
		t.Parallel()

		// Buat custom claims dengan expiry masa lalu
		now := time.Now()
		pastTime := now.Add(-1 * time.Hour)

		claims := &CustomClaims{
			UserID:  "user-123",
			IsAdmin: false,
			Type:    "access",
		}

		// Manually set ExpiresAt ke masa lalu
		claims.ExpiresAt = jwt.NewNumericDate(pastTime)

		remaining := GetRemainingTime(claims)
		assert.Equal(t, int64(0), remaining)
	})

	t.Run("Sukses - token tanpa expiry mengembalikan 0", func(t *testing.T) {
		t.Parallel()

		claims := &CustomClaims{
			UserID:  "user-123",
			IsAdmin: false,
			Type:    "access",
		}

		remaining := GetRemainingTime(claims)
		assert.Equal(t, int64(0), remaining)
	})
}
