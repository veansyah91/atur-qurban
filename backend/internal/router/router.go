package router

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/handler"
	"github.com/username/qurban-app/internal/middleware"
	"github.com/username/qurban-app/internal/repository"
	"github.com/username/qurban-app/internal/service"
	"github.com/username/qurban-app/internal/worker"
	"github.com/username/qurban-app/pkg/queue"
	"github.com/username/qurban-app/pkg/utils"
	"github.com/username/qurban-app/pkg/whatsapp"
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
	tenantRepo := repository.NewTenantRepository(db)

	// Initialize GOWA WhatsApp client
	waClient := whatsapp.NewClient(
		cfg.WhatsApp.URL,
		cfg.WhatsApp.Username,
		cfg.WhatsApp.Password,
	)

	// Initialize notification queue
	notifQueue := queue.NewNotificationQueue(rdb)

	// Initialize logger
	logger := log.New(log.Writer(), "[Server] ", log.LstdFlags)

	// Initialize services
	notificationService := service.NewNotificationService(notifQueue)
	authService := service.NewAuthService(userRepo, rdb, cfg, notificationService, logger)
	tenantService := service.NewTenantService(tenantRepo, userRepo, notificationService, logger)
	notifWorker := worker.NewNotificationWorker(notifQueue, waClient, logger)

	// Start notification worker
	notifWorker.Start(context.Background())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	tenantHandler := handler.NewTenantHandler(tenantService)

	// Grup API versi 1
	api := app.Group("/api/v1")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	authGroup.Post("/register/request-otp", authHandler.RequestRegisterOTP)
	authGroup.Post("/register/resend-otp", authHandler.ResendRegisterOTP)
	authGroup.Post("/register/verify", authHandler.VerifyRegisterOTP)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.RefreshToken)
	authGroup.Post("/forgot-password", authHandler.ForgotPassword)
	authGroup.Post("/reset-password", authHandler.ResetPassword)

	// Auth routes (protected)
	authGroup.Post("/logout", middleware.JWTAuth(authService), authHandler.Logout)
	authGroup.Get("/me", middleware.JWTAuth(authService), authHandler.Me)

	// Tenant routes (protected)
	api.Post("/tenants", middleware.JWTAuth(authService), tenantHandler.CreateTenant)
	api.Get("/tenants", middleware.JWTAuth(authService), tenantHandler.GetTenants)

	// Tenant detail routes (require membership + admin)
	tenantGroup := api.Group("/tenants/:id")
	tenantGroup.Use(middleware.JWTAuth(authService), middleware.TenantMember(tenantRepo))
	tenantGroup.Get("", tenantHandler.GetTenant)
	tenantGroup.Get("/members", tenantHandler.GetMembers)
	tenantGroup.Put("", middleware.TenantAdmin(tenantRepo), tenantHandler.UpdateTenant)
	tenantGroup.Delete("", middleware.TenantAdmin(tenantRepo), tenantHandler.DeleteTenant)
	tenantGroup.Post("/members", middleware.TenantAdmin(tenantRepo), tenantHandler.InviteMember)
	tenantGroup.Delete("/members/:user_id", middleware.TenantAdmin(tenantRepo), tenantHandler.RemoveMember)
}

