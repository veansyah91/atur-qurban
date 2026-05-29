package main

import (
"fmt"
"log"
"math/rand"

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

if err := db.AutoMigrate(&model.User{}); err != nil {
log.Fatalf("gagal migrate: %v", err)
}

if err := seedAdminDefault(db); err != nil {
log.Fatalf("gagal seed admin: %v", err)
}

if err := seedUsers(db, 50); err != nil {
log.Fatalf("gagal seed users: %v", err)
}

log.Println("seed berhasil!")
}

// hashPassword mem-hash password menggunakan bcrypt
func hashPassword(plain string) (string, error) {
hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
if err != nil {
return "", err
}
return string(hashed), nil
}

// seedAdminDefault membuat atau update user admin default
func seedAdminDefault(db *gorm.DB) error {
adminPhone := "6285271766661"
adminPassword := "9968siskom"

hashedPassword, err := hashPassword(adminPassword)
if err != nil {
return fmt.Errorf("seedAdminDefault: gagal hash password: %w", err)
}

var existingAdmin model.User
checkErr := db.Where("phone = ?", adminPhone).First(&existingAdmin).Error

if checkErr == gorm.ErrRecordNotFound {
// Buat admin baru
admin := &model.User{
Name:     "ADMIN",
Phone:    adminPhone,
IsAdmin:  true,
Password: &hashedPassword,
}
if err := db.Create(admin).Error; err != nil {
return fmt.Errorf("seedAdminDefault: create error: %w", err)
}
log.Println("admin default berhasil dibuat")
} else if checkErr == nil {
// Update password admin yang sudah ada
if err := db.Model(&existingAdmin).Update("password", hashedPassword).Error; err != nil {
return fmt.Errorf("seedAdminDefault: update password error: %w", err)
}
log.Println("password admin berhasil diupdate")
} else {
return fmt.Errorf("seedAdminDefault: check error: %w", checkErr)
}

return nil
}

// seedUsers membuat sample users dengan nama Indonesia dan password default
func seedUsers(db *gorm.DB, count int) error {
firstNames := []string{
"Ahmad", "Budi", "Citra", "Dina", "Eka", "Farah", "Gatot", "Hani", "Irfan", "Joko",
"Kamil", "Lina", "Mirna", "Nina", "Oscar", "Putri", "Qorri", "Rini", "Siti", "Tomo",
"Usman", "Vina", "Wani", "Xena", "Yani", "Zainal",
}
lastNames := []string{
"Abdulloh", "Budiman", "Chandra", "Darmawan", "Effendi", "Fauzi", "Ginting", "Handoko", "Ibnu", "Jaya",
"Kusuma", "Laksana", "Mahendra", "Nugraha", "Osman", "Priatna", "Qudus", "Rahman", "Santoso", "Taufik",
"Utama", "Vardi", "Wibisono", "Xanthis", "Yahya", "Zaini",
}

// Hash password default sekali untuk efisiensi
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

log.Printf("%d users berhasil dibuat (password: password123)", count)
return nil
}
