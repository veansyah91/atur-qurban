# Planning: Modifikasi Authentication (Security Improvement)

## Latar Belakang
Berdasarkan hasil code review, ditemukan celah keamanan kritis pada penyimpanan token di sisi frontend. Saat ini token disimpan pada `localStorage` dan dikelola manual sebagai cookie di client-side, yang mana rentan terhadap serangan XSS (Cross-Site Scripting).

## Tujuan
Meningkatkan keamanan mekanisme autentikasi dengan mencegah eksposur token ke sisi client-side JavaScript. Implementasi difokuskan pada pengiriman dan pengelolaan token melalui cookie yang aman dari sisi backend.

## Tasks High-Level (Backend)

### 1. Modifikasi Endpoint Login / Generate Token
- Ubah respons pada endpoint login/otentikasi agar tidak lagi mengembalikan `refreshToken` (dan `accessToken`) ke dalam body response JSON.
- Kirimkan token-token tersebut kepada klien menggunakan **HTTP Set-Cookie**.
- Terapkan konfigurasi cookie yang ketat:
  - **HttpOnly**: Wajib diaktifkan untuk mencegah akses baca/tulis melalui JavaScript (XSS mitigation).
  - **Secure**: Wajib diaktifkan untuk memastikan cookie hanya dikirim melalui koneksi HTTPS.
  - **SameSite**: Set ke `Lax` atau `Strict` untuk mencegah serangan CSRF.

### 2. Modifikasi Endpoint Refresh Token
- Sesuaikan endpoint pembaruan token (refresh token) agar tidak lagi menerima `refreshToken` melalui JSON body atau header otorisasi.
- Endpoint harus membaca `refreshToken` secara langsung dari **HTTP Cookie** request.
- Setelah berhasil diperbarui, token yang baru dikirimkan kembali melalui HTTP Set-Cookie sesuai dengan aturan keamanan di atas.

### 3. Modifikasi Endpoint Logout
- Sesuaikan respons endpoint logout untuk **menghapus (clear)** cookie token di sisi klien.
- Ini dapat dilakukan dengan mengirimkan header HTTP Set-Cookie untuk `accessToken` dan `refreshToken` dengan nilai kosong (`""`) dan instruksi *expired* di masa lalu.

### 4. Konfigurasi CORS & Kredensial
- Pastikan konfigurasi server / middleware CORS (Cross-Origin Resource Sharing) dikonfigurasi untuk mengizinkan kredensial (`Access-Control-Allow-Credentials: true`).
- Atur spesifik `Allowed Origins` (tidak menggunakan wildcard `*`) agar cookie dapat dikirimkan dan diterima dengan aman antara domain frontend dan backend.

## Panduan Eksekusi AI
Fokuskan modifikasi pada:
1. Handler otentikasi (login, refresh token, logout).
2. Konfigurasi server HTTP/CORS terkait keamanan kredensial.
3. Utilitas atau helper pembuatan/pengaturan cookie (Set-Cookie).