package main

import (
"fmt"
"log"
"math/rand"
"time"

"github.com/google/uuid"
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

if err := seedContacts(db); err != nil {
log.Fatalf("gagal seed contacts: %v", err)
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

now := time.Now()
var existingAdmin model.User
checkErr := db.Where("phone = ?", adminPhone).First(&existingAdmin).Error

if checkErr == gorm.ErrRecordNotFound {
// Buat admin baru dengan verified_at
admin := &model.User{
Name:       "ADMIN",
Phone:      adminPhone,
IsAdmin:    true,
Password:   &hashedPassword,
VerifiedAt: &now,
}
if err := db.Create(admin).Error; err != nil {
return fmt.Errorf("seedAdminDefault: create error: %w", err)
}
log.Println("admin default berhasil dibuat")
} else if checkErr == nil {
// Update password dan verified_at admin yang sudah ada
if err := db.Model(&existingAdmin).Updates(map[string]interface{}{
"password":    hashedPassword,
"verified_at": now,
}).Error; err != nil {
return fmt.Errorf("seedAdminDefault: update password error: %w", err)
}
log.Println("password dan verified_at admin berhasil diupdate")
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

now := time.Now()
var users []*model.User
for i := 1; i <= count; i++ {
firstName := firstNames[rand.Intn(len(firstNames))]
lastName := lastNames[rand.Intn(len(lastNames))]
phone := fmt.Sprintf("628%010d", rand.Intn(10000000000))

pwd := defaultHashedPassword
user := &model.User{
Name:       fmt.Sprintf("%s %s", firstName, lastName),
Phone:      phone,
IsAdmin:    false,
Password:   &pwd,
VerifiedAt: &now,
}
users = append(users, user)
}

if err := db.CreateInBatches(users, 10).Error; err != nil {
return fmt.Errorf("seedUsers: create error: %w", err)
}

log.Printf("%d users berhasil dibuat (password: password123, verified_at: sekarang)", count)
return nil
}

// seedContacts membuat sample contacts untuk tenant pertama yang ditemukan
func seedContacts(db *gorm.DB) error {
// Query tenant pertama yang tersedia
var tenant model.Tenant
if err := db.First(&tenant).Error; err != nil {
if err == gorm.ErrRecordNotFound {
log.Println("tidak ada tenant tersedia, skip seed contacts")
return nil
}
return fmt.Errorf("seedContacts: gagal query tenant: %w", err)
}

// Query contact category pertama untuk tenant tersebut
var category model.ContactCategory
if err := db.Where("tenant_id = ?", tenant.ID).First(&category).Error; err != nil {
if err == gorm.ErrRecordNotFound {
log.Printf("tidak ada contact category untuk tenant %s, skip seed contacts\n", tenant.Name)
return nil
}
return fmt.Errorf("seedContacts: gagal query kategori: %w", err)
}

firstNames := []string{
"Ahmad", "Budi", "Citra", "Dina", "Eka", "Farah", "Gatot", "Hani", "Irfan", "Joko",
"Kamil", "Lina", "Mirna", "Nina", "Oscar", "Putri", "Rini", "Siti", "Tomo", "Usman",
}
lastNames := []string{
"Abdulloh", "Budiman", "Chandra", "Darmawan", "Effendi", "Fauzi", "Handoko", "Ibnu", "Jaya",
"Kusuma", "Laksana", "Mahendra", "Nugraha", "Priatna", "Rahman", "Santoso", "Taufik", "Utama", "Yahya", "Zaini",
}

addresses := []string{
"Jl. Sudirman No. 10, Jakarta",
"Jl. Merdeka No. 5, Bandung",
"Jl. Imam Bonjol No. 20, Semarang",
"Jl. Diponegoro No. 8, Yogyakarta",
"Jl. Ahmad Yani No. 15, Surabaya",
}

var contacts []*model.Contact
for i := 0; i < 20; i++ {
firstName := firstNames[rand.Intn(len(firstNames))]
lastName := lastNames[rand.Intn(len(lastNames))]
name := fmt.Sprintf("%s %s", firstName, lastName)
phone := fmt.Sprintf("+628%09d", rand.Intn(1000000000))

var address *string
// Sebagian kontak memiliki alamat, sebagian tidak
if rand.Intn(2) == 0 {
addr := addresses[rand.Intn(len(addresses))]
address = &addr
}

contact := &model.Contact{
ID:                uuid.New().String(),
TenantID:          tenant.ID,
ContactCategoryID: category.ID,
Name:              name,
IsActive:          true,
Address:           address,
Phone:             &phone,
}
contacts = append(contacts, contact)
}

if err := db.CreateInBatches(contacts, 10).Error; err != nil {
return fmt.Errorf("seedContacts: create error: %w", err)
}

log.Printf("20 contacts berhasil dibuat untuk tenant '%s' (kategori: %s)\n", tenant.Name, category.Name)
return nil
}
