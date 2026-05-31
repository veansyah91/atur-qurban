# Planning: Perbaikan Halaman Home 404

## Deskripsi Masalah
Saat mengakses halaman utama (`http://localhost:3000/`), aplikasi memunculkan halaman 404 (Not Found), padahal file untuk halaman utama (`src/app/page.tsx`) sudah ada.

## Analisis High Level
Dalam proyek berbasis Next.js App Router, rute `/` di-*handle* oleh `src/app/page.tsx`. Tampilan 404 biasanya terjadi akibat beberapa hal:

1. **Cache / Build Issues:** Terdapat status build atau cache routing yang rusak/kadaluwarsa di dalam direktori `.next`.
2. **Gangguan Middleware:** Konfigurasi regex pada `matcher` di `src/middleware.ts` mungkin secara tak terduga memblokir atau me-redirect request menuju `/`.
3. **Kesalahan Layout Utama:** File `src/app/layout.tsx` gagal me-render props `children` (komponen halaman utama) dengan benar, memicu proses fallback internal Next.js ke halaman 404.
4. **Struktur File/Konflik:** Adanya folder tak wajar, atau *conflict* pada route (seperti file `route.ts` yang berdampingan dengan `page.tsx` di root `app`).

## Instruksi Implementasi

Kepada AI Pelaksana, jalankan langkah-langkah *troubleshooting* dan perbaikan berikut:

### 1. Bersihkan Cache dan Jalankan Ulang
- Hapus folder `.next` sepenuhnya.
- Jalankan ulang *development server* (contoh: `npm run dev`).
- Uji kembali rute `/`. Jika berhasil me-render UI, proses selesai.

### 2. Evaluasi Middleware
- Coba matikan middleware dengan merename file `src/middleware.ts` menjadi nama lain (misal `middleware.ts.bak`).
- Cek kembali rute `/`. Jika berhasil dimuat, berarti masalahnya ada di `middleware.ts`.
- *Fix*: Perbaiki pola *regex* di dalam konfigurasi `matcher` agar me-whitelist atau tidak menangkap request menuju index `/`.

### 3. Validasi Root Layout & Ekspor Komponen
- Pastikan komponen pada `src/app/page.tsx` diekspor secara valid (`export default function...`).
- Cek kode pada `src/app/layout.tsx`. Pastikan ia merender parameter `{ children }` di dalam tag HTML/Body, sehingga konten dari `page.tsx` benar-benar disuntikkan ke halaman.
- Pastikan tidak ada konflik file (seperti adanya `route.ts` di direktori `src/app/`).

## Kriteria Sukses (Acceptance Criteria)
- Mengakses URL `http://localhost:3000/` harus me-return HTTP 200 OK dan me-render *Landing Page* sesuai kode yang ada di `src/app/page.tsx`.
- Log *development server* bersih dari error saat me-render halaman utama.