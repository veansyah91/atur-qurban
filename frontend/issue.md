# Issue Planning: Fix Authentication API Endpoints (404 Errors)

## Deskripsi Masalah
Ditemukan error 404 (Not Found) pada dua fitur utama autentikasi:
1. **Login**: Request POST ke `/auth/login` mengembalikan 404.
2. **Register**: Request POST ke `/auth/register/request-otp` mengembalikan 404.

Hal ini mengindikasikan ketidaksesuaian antara endpoint yang dipanggil di frontend dengan struktur API yang tersedia di backend.

## Rencana Perbaikan

### 1. Sinkronisasi Endpoint Login
*   **File**: `src/services/auth.service.ts`
*   **Instruksi**: Periksa dan sesuaikan path endpoint pada fungsi `loginAdmin`. 
*   **Analisis**: Endpoint `/auth/login` kemungkinan besar salah. Berdasarkan pola service lainnya, periksa apakah seharusnya menggunakan `/auth/admin/login` atau cukup `/login`. Sesuaikan dengan dokumentasi API backend yang berlaku.

### 2. Sinkronisasi Endpoint Register OTP
*   **File**: `src/services/auth-register.service.ts`
*   **Instruksi**: Periksa dan sesuaikan path endpoint pada fungsi `requestRegisterOtp`.
*   **Analisis**: Endpoint `/auth/register/request-otp` menyebabkan 404. Perhatikan konsistensi penamaan dengan fitur OTP lainnya (misal: di `auth.service.ts` menggunakan pola `/auth/otp/...`). Ubah path menjadi endpoint yang valid sesuai kontrak API backend (contoh: `/auth/otp/register`).

### 3. Verifikasi Konfigurasi Base URL
*   **File**: `.env` (atau environment variable)
*   **Instruksi**: Pastikan `NEXT_PUBLIC_API_URL` mengarah ke host backend yang benar (bukan ke server Next.js itu sendiri jika API dilayani oleh service terpisah).

## Kriteria Keberhasilan
- Tidak ada lagi error 404 saat melakukan submit login.
- Request OTP untuk registrasi berhasil terkirim ke backend dengan status 200/201.
- Fungsi-fungsi terkait (seperti `verifyOtp`) tetap bekerja dengan normal.