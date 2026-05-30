# Copilot Instructions — Atur Qurban Backend

## Konteks Project
Backend REST API untuk aplikasi manajemen qurban. Ditulis dalam Go, menggunakan arsitektur berlapis (layered architecture): handler → service → repository.

## Tech Stack
| Komponen | Library / Versi |
|---|---|
| Language | Go 1.22 |
| HTTP Framework | Fiber v2 (`github.com/gofiber/fiber/v2`) |
| ORM | GORM v1.25 (`gorm.io/gorm`) |
| Database | PostgreSQL 16 (`gorm.io/driver/postgres`) |
| Cache | Redis 7 (`github.com/redis/go-redis/v9`) |
| Config | Viper v1.19 (`github.com/spf13/viper`) |
| Migration | golang-migrate (via CLI `migrate`) |
| Live Reload | Air (`air`) |
| Container | Docker Compose |

## Struktur Folder
```
backend/
├── cmd/
│   └── main.go              # Entry point — inisialisasi app, DB, Redis, router
├── config/
│   └── config.go            # Struct & LoadConfig() menggunakan Viper
├── internal/
│   ├── handler/             # HTTP handler (Fiber), 1 file per domain
│   ├── middleware/          # Middleware Fiber (auth, logging, dll)
│   ├── model/               # GORM model / entitas domain
│   ├── repository/          # Query database (interface + implementasi)
│   ├── router/
│   │   └── router.go        # Registrasi semua route
│   └── service/             # Business logic (interface + implementasi)
├── migrations/              # File SQL migrasi (up/down), dikelola golang-migrate
├── pkg/
│   ├── database/
│   │   ├── postgres.go      # ConnectPostgres() — koneksi GORM + connection pool
│   │   └── redis.go         # ConnectRedis() — koneksi dengan Ping verify
│   └── utils/
│       └── response.go      # SuccessResponse, ErrorResponse, PaginateResponse
├── docs/                    # Output Swagger (generate via `make swag`)
├── .env                     # Konfigurasi lokal (tidak di-commit)
├── .env.example             # Template konfigurasi
├── .air.toml                # Konfigurasi live reload Air
├── Dockerfile               # Multi-stage build (builder + alpine runner)
├── docker-compose.yml       # PostgreSQL 16 + Redis 7 + app service
└── Makefile                 # Target: run, build, test, tidy, migrate-*, docker-*, swag
```

## Module Name
```
github.com/username/qurban-app
```

## Coding Conventions

### Bahasa
- **Semua komentar kode dalam bahasa Indonesia**
- Nama variabel, fungsi, dan tipe tetap dalam bahasa Inggris

### Response Format
Selalu gunakan helper dari `pkg/utils/response.go`:
```go
// Sukses
utils.SuccessResponse(c, "pesan sukses", data)

// Error
utils.ErrorResponse(c, fiber.StatusBadRequest, "pesan error")

// Paginated
utils.PaginateResponse(c, "pesan", data, meta)
```

### Primary Key
- Semua model menggunakan **UUID** sebagai primary key
- Gunakan `github.com/google/uuid` (sudah ada di go.sum sebagai indirect dep)

### Database
- **Tidak ada global variable untuk DB atau Redis** — selalu inject via parameter/constructor
- Konfigurasi connection pool sudah diatur di `ConnectPostgres()`: max 25 open, 5 idle, 5 menit lifetime
- Redis diverifikasi dengan Ping saat inisialisasi

### Dependency Injection
Alur injeksi dari `main.go`:
```
main.go → ConnectPostgres/ConnectRedis → NewRouter(app, db, rdb) → handler(service(repo(db)))
```

### Error Handling
- Kembalikan error dengan konteks: `fmt.Errorf("FunctionName: %w", err)`
- Gunakan `log.Fatalf` di main.go untuk error fatal saat startup

### Router
- API versioning: semua endpoint bisnis di bawah prefix `/api/v1`
- Health check di `/health` (cek koneksi DB + Redis)

### Makefile Commands
| Command | Fungsi |
|---|---|
| `make run` | Jalankan server dengan Air (live reload, port 8080) |
| `make build` | Build binary ke `./bin/app` |
| `make test` | Jalankan test dengan verbose + coverage |
| `make tidy` | `go mod tidy` |
| `make docker-up` | `docker compose up -d` |
| `make docker-down` | `docker compose down` |
| `make migrate-up` | Terapkan semua migration |
| `make migrate-down` | Batalkan 1 migration terakhir |
| `make migrate-create name=xxx` | Buat file migration baru |
| `make swag` | Generate dokumentasi Swagger ke `./docs` |

### Environment Variables
Konfigurasi wajib di `.env` (lihat `.env.example`):
- `APP_ENV`, `APP_PORT` (default 8080), `APP_SECRET`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `JWT_SECRET`, `JWT_EXPIRY_HOUR`, `JWT_REFRESH_EXPIRY_DAY`

---

## Planning

Sebelum mengimplementasikan fitur baru, buat dokumen perencanaan di session folder (`plan.md`). Dokumen ini akan digunakan oleh AI model lain untuk implementasi.

### Format Plan

- **Ringkas & actionable** — gunakan bullet point, bukan paragraf panjang
- **Daftar file** yang perlu dibuat atau dimodifikasi
- **Urutan implementasi** jika ada dependensi antar komponen
- **Catat keputusan penting** (misal: apakah perlu migration baru, middleware baru, dll)

### Isi Minimal Plan

1. Deskripsi singkat fitur
2. File baru yang perlu dibuat (model, repo, service, handler, migration)
3. File yang perlu dimodifikasi (router, dll)
4. Endpoint dan method HTTP
5. Catatan khusus jika ada (validasi, auth, dsb)

---

## Unit Testing

### Prinsip Utama

- **Tidak ada koneksi nyata** ke database, Redis, atau layanan eksternal
- Gunakan **mock** untuk semua dependensi eksternal
- Setiap layer ditest secara terpisah (service ≠ handler ≠ repository)

### Library

- `testing` — bawaan Go
- `github.com/stretchr/testify` — assertions (`assert`, `require`) dan mock (`mock`)

### Struktur Mock

Buat mock dari interface repository/service menggunakan `testify/mock`:

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(id string) (*model.User, error) {
    args := m.Called(id)
    return args.Get(0).(*model.User), args.Error(1)
}
```

### Testing Service Layer

- Inject mock repository ke service
- Test business logic dan error handling
- Gunakan table-driven tests untuk banyak skenario

```go
func TestUserService_FindByID(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        mock    func(*MockUserRepository)
        wantErr bool
    }{
        // ... test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockUserRepository)
            tt.mock(repo)
            svc := NewUserService(repo)
            _, err := svc.FindByID(tt.input)
            assert.Equal(t, tt.wantErr, err != nil)
        })
    }
}
```

### Testing Handler Layer

- Gunakan `net/http/httptest` atau Fiber test helper
- Inject mock service ke handler
- Verifikasi HTTP status code dan body response

```go
func TestUserHandler_GetByID(t *testing.T) {
    app := fiber.New()
    mockSvc := new(MockUserService)
    mockSvc.On("FindByID", "123").Return(&model.User{}, nil)

    h := NewUserHandler(mockSvc)
    app.Get("/users/:id", h.GetByID)

    req := httptest.NewRequest("GET", "/users/123", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
```

### Konvensi File Test

- File test berada di package yang sama: `user_service_test.go`
- Nama fungsi test: `TestNamaStruct_NamaMethod`
- Jalankan dengan: `make test`

---

## Unit Testing — Tambahan Konvensi

### Target coverage
- Minimal 80% per file di `internal/service/`
- Jalankan: `go test ./internal/service/... -cover` untuk cek

### Lokasi file mock
- Semua mock disimpan di `internal/mocks/`
- Nama file: `mock_nama_interface.go`
  (contoh: `mock_user_repository.go`, `mock_auth_service.go`)

### Skenario wajib per method
Setiap service method MINIMAL harus punya test untuk:
- **Happy path** — input valid, semua dependency sukses
- **Input tidak valid** — format salah, field kosong
- **Resource tidak ditemukan** — dependency return not found error
- **Dependency error** — dependency return error umum

---

## Code Review Checkpoint

Jalankan review ini setelah setiap fitur selesai, sebelum lanjut ke fitur berikutnya:

### Checklist per fitur
- [ ] Semua constraint di `copilot-instructions.md` diikuti
- [ ] Tidak ada business logic di handler
- [ ] Semua error di-wrap dengan `fmt.Errorf("FunctionName: %w", err)`
- [ ] Tidak ada data sensitif (password, token) yang masuk ke log
- [ ] Semua endpoint yang butuh auth sudah pasang middleware JWT
- [ ] Response selalu menggunakan helper dari `pkg/utils/response.go`
- [ ] Setiap handler punya Swagger annotation lengkap

### Prompt review standar
Gunakan prompt ini setelah setiap fitur selesai:

> "@workspace Review [nama file handler] dan [nama file service] yang baru dibuat. Periksa apakah semua konvensi di copilot-instructions.md sudah diikuti. Tampilkan temuan beserta baris yang bermasalah dan saran perbaikannya."

---

## Auth Convention

### Model User — hal penting
- `ID` bertipe string (UUID disimpan sebagai string, bukan `uuid.UUID`)
- `Password` nullable (`*string`) — user tanpa password valid untuk flow OTP
- `Email` nullable (`*string`) — opsional saat registrasi
- Tidak ada field `role` — gunakan `IsAdmin` bool untuk membedakan hak akses
- Nonaktif user dilakukan dengan soft delete (`DeletedAt`), bukan toggle field
- `LastLoginAt` wajib diupdate setiap login sukses

### Identitas login
- Primary identifier: `Phone` dalam format E.164
  (contoh: `+6281234567890`, bukan `081234567890`)
- Validasi format menggunakan library `nyaruka/phonenumbers`

### Strategi auth berdasarkan IsAdmin
- `IsAdmin = true` → admin, login dengan `Phone` + `Password` (wajib ada)
- `IsAdmin = false` → peserta, login dengan `Phone` + OTP
- Password nil pada user `IsAdmin=true` adalah kondisi invalid — tolak login

### JWT payload standar
```json
{
  "user_id":  "string (uuid)",
  "phone":    "+62...",
  "is_admin": true | false,
  "exp":      "unix_timestamp"
}
```

### Middleware auth
- `JWTMiddleware` — validasi token, set locals `user_id` & `is_admin`
- `RequireAdmin` — cek locals `is_admin == true`, return 403 jika bukan admin
- Jangan gunakan role string — selalu gunakan `is_admin` dari JWT locals

### OTP
- Panjang: 6 digit angka
- Expired: 5 menit sejak dibuat
- Maksimal percobaan: 3x
- Simpan hash OTP ke tabel `otp_codes` (bukan nilai asli)

### Setelah login sukses
- Update field `LastLoginAt` ke waktu sekarang
- Simpan refresh token hash ke tabel `refresh_tokens`

---

## Swagger Annotation

Setiap handler wajib punya Swagger annotation lengkap:

```go
// @Summary     Ringkasan singkat endpoint
// @Description Deskripsi lebih detail
// @Tags        nama-group
// @Accept      json
// @Produce     json
// @Param       body body NamaRequestStruct true "Deskripsi body"
// @Success     200 {object} utils.BaseResponse
// @Failure     400 {object} utils.BaseResponse
// @Failure     401 {object} utils.BaseResponse
// @Router      /api/v1/path [method]
```
