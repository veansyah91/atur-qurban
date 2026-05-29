package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/router"
	"github.com/username/qurban-app/pkg/database"
)

func main() {
	// Muat konfigurasi dari .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("gagal memuat konfigurasi: %v", err)
	}

	// Hubungkan ke PostgreSQL
	db, err := database.ConnectPostgres(&cfg.Database, cfg.App.Env)
	if err != nil {
		log.Fatalf("gagal terhubung ke PostgreSQL: %v", err)
	}

	// Hubungkan ke Redis
	rdb, err := database.ConnectRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("gagal terhubung ke Redis: %v", err)
	}
	defer rdb.Close()

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "Qurban App v1.0",
	})

	// Daftarkan semua route
	router.NewRouter(app, db, rdb)

	// Jalankan server di goroutine terpisah agar graceful shutdown bisa bekerja
	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("server berhenti: %v", err)
		}
	}()

	log.Printf("server berjalan di port %s (env: %s)", cfg.App.Port, cfg.App.Env)

	// Tunggu sinyal SIGINT atau SIGTERM untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("mematikan server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("gagal mematikan server: %v", err)
	}

	log.Println("server berhasil dimatikan")
}
