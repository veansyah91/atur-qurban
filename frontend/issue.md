# Planning Fitur Authentication

Dokumen ini berisi high-level planning untuk implementasi fitur authentication pada aplikasi. Referensi API yang digunakan berdasarkan dokumentasi backend (`backend-documentation/auth.md`).

## 1. Persiapan UI dan Aset
- Gunakan asset-asset image yang tersedia di repositori (seperti logo, background, ikon, dll) untuk kebutuhan halaman authentication agar UI terlihat menarik dan representatif.
- Buat struktur komponen UI untuk:
  - Halaman Login
  - Halaman Registrasi
  - Halaman Verifikasi OTP
  - Halaman Lupa Password (dan form input OTP + Password Baru)

## 2. Alur (Flow) Authentication

### A. Alur Registrasi
Sesuai instruksi, flow register harus berjalan sebagai berikut: **Halaman Register -> Halaman Verifikasi**.
1. **Request OTP (Register)**: Pengguna mengisi Form Registrasi (nama, nomor telepon, password, confirm password). Jika sukses, pengguna langsung diarahkan ke **Halaman Verifikasi**.
2. **Kirim Ulang (Resend) OTP**: Di Halaman Verifikasi, sediakan tombol kirim ulang OTP yang bisa diklik (biasanya ditambahkan countdown) untuk memanggil API resend OTP.
3. **Verify OTP**: Pengguna memasukkan OTP. Jika sukses, simpan token autentikasi dan arahkan masuk ke halaman utama aplikasi.

### B. Alur Login & Sesi
1. **Login**: Pengguna memasukkan nomor telepon dan password. Jika sukses, simpan token dan masuk ke aplikasi.
2. **Refresh Token**: Implementasikan mekanisme pembaruan token (refresh token). Jika *access token* kedaluwarsa, secara otomatis request token baru agar sesi pengguna tidak terputus.
3. **Get Current User (Me)**: Pastikan memanggil API profil (Me) saat aplikasi di-refresh untuk validasi sesi dan state management pengguna.
4. **Logout**: Bersihkan token dan kembali ke halaman Login.
5. **Lupa Password**: Sediakan flow untuk request reset password (Kirim OTP) dilanjutkan dengan mengubah password baru menggunakan OTP.

## 3. Optimasi dan Praktik Terbaik
- **Penyimpanan Token**: Simpan `access_token` dan `refresh_token` secara aman.
- **API Interceptor**: Implementasikan HTTP interceptor (misal pada Axios/Fetch API) untuk otomatis menyematkan header `Authorization: Bearer <ACCESS_TOKEN>`.
- **Handling Error & Validasi**: Pastikan pesan error dari API ditangani dengan baik dan ditampilkan ke pengguna. Normalisasi nomor telepon dari input pengguna agar kompatibel dengan sistem.
- **State Management**: Gunakan global state/context untuk mengelola status autentikasi pengguna dan mengamankan rute (Protected Routes).

---
**Instruksi untuk AI Implementer:**
Implementasikan planning ini secara komprehensif. Mulai dari setup routing, pembuatan komponen UI dengan aset yang ada, integrasi API service, hingga manajemen state autentikasi dengan mempertimbangkan best practices yang disebutkan di atas.