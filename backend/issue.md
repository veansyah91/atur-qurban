# Issue: Modifikasi Tenant — Upload Logo

## Deskripsi
Modifikasi fitur tenant agar logo dikirim sebagai file upload (bukan URL string), disimpan di local storage atau S3, dengan validasi ukuran maksimal 1MB.

---

## Perubahan yang Dibutuhkan

### 1. Storage Abstraction (`pkg/storage/`)

Buat package baru untuk abstraksi penyimpanan file:

- **`storage.go`** — interface `StorageService` dengan method `Upload(file multipart.File, filename, contentType string) (string, error)` dan `Delete(path string) error`. Sertakan factory `NewStorageService(cfg StorageConfig)` yang memilih implementasi berdasarkan `STORAGE_DRIVER`.
- **`local.go`** — implementasi local disk: simpan file ke folder `STORAGE_LOCAL_PATH` (default: `./uploads`), kembalikan URL publik berdasarkan base URL app.
- **`s3.go`** — implementasi AWS S3 menggunakan `aws-sdk-go-v2`: upload ke bucket `S3_BUCKET`, kembalikan URL publik.

### 2. Konfigurasi (`config/config.go` + `.env.example`)

Tambah `StorageConfig` ke `AppConfig`:

```
STORAGE_DRIVER=local         # "local" atau "s3"
STORAGE_LOCAL_PATH=./uploads  # path folder untuk local storage
STORAGE_BASE_URL=http://localhost:8080/uploads  # base URL untuk serve file lokal

# Khusus S3 (opsional, hanya jika STORAGE_DRIVER=s3)
S3_BUCKET=
S3_REGION=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_ENDPOINT=   # opsional, untuk S3-compatible (misal MinIO)
```

### 3. Serve Static Files (`router/router.go`)

Tambahkan `app.Static("/uploads", "./uploads")` agar file lokal bisa diakses via HTTP.

### 4. Handler Tenant (`internal/handler/tenant_handler.go`)

Ubah `CreateTenant` dan `UpdateTenant`:

- Ganti `Accept: json` → `Accept: multipart/form-data`
- Field `name`, `address`, `description` dikirim sebagai form field
- Field `logo` dikirim sebagai file upload (`multipart/form-file`)
- Validasi: ukuran file ≤ 1MB, MIME type harus `image/*`
- Jika ada logo lama saat update, hapus file lama via `StorageService.Delete()`
- Inject `StorageService` ke `TenantHandler` via constructor

### 5. Router (`internal/router/router.go`)

- Inisialisasi `StorageService` dari config
- Pass ke `handler.NewTenantHandler(tenantService, storageService)`

### 6. Testing (`internal/handler/tenant_handler_test.go`)

Buat test untuk `CreateTenant` dan `UpdateTenant` di handler layer:

- Mock `TenantService` dan `StorageService`
- Skenario wajib:
  - Upload logo valid → sukses
  - File melebihi 1MB → 400
  - MIME type bukan image → 400
  - Tanpa file logo → sukses (logo opsional)
  - StorageService error → 500
- Gunakan `multipart.Writer` untuk build request body di test
- Tidak ada koneksi ke DB atau storage nyata

### 7. Dokumentasi (`documentation/tenant.md`)

Update endpoint `POST /tenants` dan `PUT /tenants/:id`:

- Ubah `Content-Type` dari `application/json` → `multipart/form-data`
- Tambahkan constraint logo: maks 1MB, format image
- Update contoh request

---

## Urutan Implementasi

1. `config/config.go` + `.env.example` — tambah StorageConfig
2. `pkg/storage/` — interface, local, s3
3. `router/router.go` — init storage, serve static
4. `internal/handler/tenant_handler.go` — ubah handler
5. `internal/handler/tenant_handler_test.go` — tulis test
6. `documentation/tenant.md` — update docs

---

## Catatan

- `TenantService` **tidak perlu diubah** — tetap menerima `*string` untuk logo (handler yang upload, lalu pass URL ke service)
- Logo bersifat opsional — create/update tanpa logo tetap valid
- Saat update, jika logo baru dikirim: upload baru, hapus lama. Jika tidak dikirim: biarkan logo lama
- `go.mod` perlu tambah dependency `aws-sdk-go-v2` jika S3 diimplementasikan
