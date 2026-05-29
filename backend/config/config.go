package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig menyimpan seluruh konfigurasi aplikasi
type AppConfig struct {
	App      AppSetting
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// AppSetting konfigurasi umum aplikasi
type AppSetting struct {
	Env    string
	Port   string
	Secret string
}

// DatabaseConfig konfigurasi koneksi PostgreSQL
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig konfigurasi koneksi Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// JWTConfig konfigurasi token JWT
type JWTConfig struct {
	Secret           string
	ExpiryHour       int
	RefreshExpiryDay int
}

// LoadConfig memuat konfigurasi dari file .env menggunakan Viper
func LoadConfig() (*AppConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}

	cfg := &AppConfig{
		App: AppSetting{
			Env:    viper.GetString("APP_ENV"),
			Port:   viper.GetString("APP_PORT"),
			Secret: viper.GetString("APP_SECRET"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret:           viper.GetString("JWT_SECRET"),
			ExpiryHour:       viper.GetInt("JWT_EXPIRY_HOUR"),
			RefreshExpiryDay: viper.GetInt("JWT_REFRESH_EXPIRY_DAY"),
		},
	}

	return cfg, nil
}
