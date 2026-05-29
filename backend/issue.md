# Issue: Auth — Modifikasi Password & Forgot Password

## Konteks
Auth feature sudah ada (register, login, refresh, logout, me). Modifikasi berikut perlu diterapkan:

---

## 1. Perubahan Model User

- Field `password` yang sebelumnya nullable → sekarang **wajib diisi** saat register.
- Password di-hash menggunakan **bcrypt** sebelum disimpan ke DB.
- Di response API, field `password` tidak pernah dikembalikan.

---

## 2. Modifikasi Register

Endpoint: `POST /api/v1/auth/register`

- Tambah field `password` pada request body (wajib).
- Hash password dengan bcrypt sebelum simpan ke DB.

---

## 3. Modifikasi Login

Endpoint: `POST /api/v1/auth/login`

- Tambah field `password` pada request body (wajib).
- Setelah user ditemukan by phone, verifikasi password dengan bcrypt compare.
- Jika tidak cocok → 401 Unauthorized.

---

## 4. Fitur Lupa Sandi (Forgot Password)

### 4a. Request OTP — `POST /api/v1/auth/forgot-password`

Request body:
```json
{ "phone": "string" }
```

Flow:
1. Normalisasi phone
2. Cek user ada by phone → 404 jika tidak ada
3. Generate OTP 6 digit
4. Simpan ke Redis: key `otp:reset:<phone>`, TTL 5 menit
5. Return OTP dalam response (di production, akan dikirim via WhatsApp/SMS)

### 4b. Reset Password — `POST /api/v1/auth/reset-password`

Request body:
```json
{ "phone": "string", "otp": "string", "new_password": "string" }
```

Flow:
1. Normalisasi phone
2. Ambil OTP dari Redis
3. Bandingkan OTP → 400 jika salah atau expired
4. Hash new_password dengan bcrypt
5. Update password user di DB
6. Hapus OTP dari Redis
7. Return 200 OK

---

## 5. Repository

Tambah method baru: `UpdatePassword(userID, hashedPassword string) error`

---

## 6. Modifikasi DB Seed

- **50 sample users**: setiap user mendapat password default yang di-hash (misalnya `password123`)
- **Admin default**: password di-update menjadi `9968siskom` (di-hash bcrypt)
- Seed bersifat idempotent: gunakan `ON CONFLICT` atau cek exist sebelum insert/update

---

## Checklist Implementasi

- [x] Modifikasi `internal/service/auth_service.go` — Register & Login + password logic + ForgotPassword + ResetPassword
- [x] Modifikasi `internal/handler/auth_handler.go` — update request struct + tambah 2 handler baru
- [x] Tambah method `UpdatePassword` di `internal/repository/user_repository.go`
- [x] Update `internal/router/router.go` — daftarkan 2 route baru
- [x] Update `cmd/seed/main.go` — tambah password di seed users & admin
