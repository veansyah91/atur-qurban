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
