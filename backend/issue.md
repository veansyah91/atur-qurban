# Issue: Perbaiki CORS Error pada Endpoint Login

## Deskripsi
Terdapat isu CORS (Cross-Origin Resource Sharing) ketika frontend mencoba mengakses endpoint login pada backend. Request diblokir karena tidak adanya header `Access-Control-Allow-Origin` pada response preflight.

**Error Message:**
> Access to XMLHttpRequest at 'http://localhost:8080/api/v1/auth/login' from origin 'http://localhost:3000' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource.

## Objective
Mengatasi masalah CORS agar aplikasi frontend (`http://localhost:3000`) dapat berhasil melakukan request ke backend (`http://localhost:8080`) tanpa diblokir oleh browser.

## Instruksi High-Level untuk Implementasi
1. **Cari Konfigurasi CORS**: Temukan file atau bagian kode di backend yang menangani konfigurasi CORS, middleware global, atau pengaturan keamanan.
2. **Izinkan Origin Frontend**: Update konfigurasi tersebut dengan menambahkan `http://localhost:3000` ke dalam daftar *Allowed Origins* (hindari penggunaan wildcard `*` di tahap produksi).
3. **Izinkan Method & Header**: Pastikan pengaturan CORS juga mengizinkan HTTP Methods yang diperlukan (khususnya `OPTIONS` untuk preflight, dan `POST` untuk login) serta HTTP Headers yang relevan (seperti `Content-Type`, `Authorization`).
4. **Verifikasi**: Pastikan bahwa error CORS hilang di browser dengan menguji request login dari frontend ke endpoint backend.