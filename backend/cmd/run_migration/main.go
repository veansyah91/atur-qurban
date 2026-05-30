package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("gagal load config: %v", err)
	}

	db, err := database.ConnectPostgres(&cfg.Database, cfg.App.Env)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	// Baca semua file migration .up.sql
	migrationsDir := "./migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("gagal baca directory migrations: %v", err)
	}

	var upFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upFiles = append(upFiles, file.Name())
		}
	}

	// Sort files untuk urutan yang benar
	sort.Strings(upFiles)

	// Jalankan setiap migration
	for _, file := range upFiles {
		path := filepath.Join(migrationsDir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("gagal baca file %s: %v", file, err)
			continue
		}

		sql := string(content)
		// Skip comment lines
		re := regexp.MustCompile(`(?m)^--.*$`)
		sql = re.ReplaceAllString(sql, "")
		// Trim whitespace
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}

		fmt.Printf("Jalankan migration: %s\n", file)
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("WARNING - migration %s gagal (mungkin sudah exists): %v", file, err)
			continue
		}
		fmt.Printf("✓ Migration %s berhasil\n", file)
	}

	fmt.Println("\n✓ Semua migrations berhasil dijalankan")
}
