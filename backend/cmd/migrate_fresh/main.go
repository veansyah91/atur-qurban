package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"

	"github.com/username/qurban-app/config"
	"github.com/username/qurban-app/internal/model"
	"github.com/username/qurban-app/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	// Drop tables (fresh)
	if err := db.Migrator().DropTable(&model.User{}); err != nil {
		log.Fatalf("gagal drop table users: %v", err)
	}
	log.Println("tables dropped")

	// Migrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("gagal migrate: %v", err)
	}
	log.Println("migrations applied")

	// Seed
	if err := seedAdminDefault(db); err != nil {
		log.Fatalf("gagal seed admin: %v", err)
	}
	if err := seedUsers(db, 50); err != nil {
		log.Fatalf("gagal seed users: %v", err)
	}

	log.Println("migrate:fresh --seed selesai")
}

func hashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func seedAdminDefault(db *gorm.DB) error {
	adminPhone := "085271766661"  // Gunakan format lokal, akan di-normalize di DB
	adminPassword := "9968siskom"

	hashedPassword, err := hashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("seedAdminDefault: gagal hash password: %w", err)
	}

	// Normalisasi phone sebelum simpan
	normalizedPhone := normalizePhone(adminPhone)
	
	admin := &model.User{
		Name:     "ADMIN",
		Phone:    normalizedPhone,
		IsAdmin:  true,
		Password: &hashedPassword,
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("seedAdminDefault: create error: %w", err)
	}
	return nil
}

func normalizePhone(phone string) string {
	// Remove non-digits
	phone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	// Convert 08... to 62...
	if len(phone) > 2 && phone[0] == '0' && phone[1] == '8' {
		phone = "62" + phone[2:]
	}
	return phone
}

func seedUsers(db *gorm.DB, count int) error {
	firstNames := []string{"Ahmad", "Budi", "Citra", "Dina", "Eka", "Farah", "Gatot", "Hani", "Irfan", "Joko",
		"Kamil", "Lina", "Mirna", "Nina", "Oscar", "Putri", "Qorri", "Rini", "Siti", "Tomo",
		"Usman", "Vina", "Wani", "Xena", "Yani", "Zainal"}
	lastNames := []string{"Abdulloh", "Budiman", "Chandra", "Darmawan", "Effendi", "Fauzi", "Ginting", "Handoko", "Ibnu", "Jaya",
		"Kusuma", "Laksana", "Mahendra", "Nugraha", "Osman", "Priatna", "Qudus", "Rahman", "Santoso", "Taufik",
		"Utama", "Vardi", "Wibisono", "Xanthis", "Yahya", "Zaini"}

	defaultHashedPassword, err := hashPassword("password123")
	if err != nil {
		return fmt.Errorf("seedUsers: gagal hash default password: %w", err)
	}

	var users []*model.User
	for i := 1; i <= count; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		phone := fmt.Sprintf("628%010d", rand.Intn(10000000000))

		pwd := defaultHashedPassword
		user := &model.User{
			Name:     fmt.Sprintf("%s %s", firstName, lastName),
			Phone:    phone,
			IsAdmin:  false,
			Password: &pwd,
		}
		users = append(users, user)
	}

	if err := db.CreateInBatches(users, 10).Error; err != nil {
		return fmt.Errorf("seedUsers: create error: %w", err)
	}
	return nil
}
