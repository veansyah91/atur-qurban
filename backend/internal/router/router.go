package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/handler"
	"github.com/username/qurban-app/internal/middleware"
	"github.com/username/qurban-app/internal/repository"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/pkg/utils"
	"gorm.io/gorm"
)

// NewRouter mendaftarkan semua route aplikasi
func NewRouter(app *fiber.App, db *gorm.DB, rdb *redis.Client) {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("gagal load config: " + err.Error())
	}

	// Health check — cek koneksi DB dan Redis
	app.Get("/health", func(c *fiber.Ctx) error {
		// Cek koneksi PostgreSQL
		sqlDB, err := db.DB()
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "database tidak tersedia")
		}
		if err := sqlDB.Ping(); err != nil {
			return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "database tidak merespons")
		}

		// Cek koneksi Redis
		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "redis tidak merespons")
		}

		return utils.SuccessResponse(c, "server berjalan normal", fiber.Map{
			"database": "ok",
			"redis":    "ok",
		})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, rdb, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Grup API versi 1
	api := app.Group("/api/v1")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.RefreshToken)
	authGroup.Post("/forgot-password", authHandler.ForgotPassword)
	authGroup.Post("/reset-password", authHandler.ResetPassword)

	// Auth routes (protected)
	authGroup.Post("/logout", middleware.JWTAuth(authService), authHandler.Logout)
	authGroup.Get("/me", middleware.JWTAuth(authService), authHandler.Me)
}

