# Tenant API Documentation

Base URL: http://localhost:8080/api/v1

Notes:
- Semua endpoint membutuhkan `Authorization: Bearer <ACCESS_TOKEN>`.
- Tenant pertama yang dibuat user → `status: free`, tidak ada tanggal expired.
- Tenant ke-2 dan seterusnya → `status: paid`, expired +14 hari dari tanggal pembuatan.
- Setelah tenant berhasil dibuat, notifikasi WhatsApp dikirim ke nomor pembuat.

---

## Endpoints

### 1) Buat Tenant
- Method: POST
- Path: /tenants
- Auth: Bearer token
- Request JSON:
```json
{ "name": "Kelompok Masjid Al-Ikhlas", "logo": "https://example.com/logo.png" }
```
- Success (201) sample:
```json
{
  "success": true,
  "message": "tenant berhasil dibuat",
  "data": {
    "id": "uuid",
    "name": "Kelompok Masjid Al-Ikhlas",
    "slug": "kelompok-masjid-al-ikhlas",
    "status": "free",
    "logo": "https://example.com/logo.png",
    "expired_at": null,
    "owner_id": "uuid",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### 2) List Tenant Milik User
- Method: GET
- Path: /tenants
- Auth: Bearer token
- Description: Mengembalikan semua tenant di mana user adalah member.
- Success (200): array tenant

### 3) Detail Tenant
- Method: GET
- Path: /tenants/:id
- Auth: Bearer token + harus member tenant
- Success (200): data tenant

### 4) Update Nama Tenant
- Method: PUT
- Path: /tenants/:id
- Auth: Bearer token + harus admin tenant
- Request JSON:
```json
{
  "name": "Nama Baru",
  "address": "Jl. Sudirman No. 1, Jakarta",
  "description": "Deskripsi kelompok qurban",
  "logo": "https://example.com/logo-baru.png"
}
```
- Catatan: `address`, `description`, dan `logo` bersifat **opsional** — boleh tidak dikirim atau di-set ke `null`.
- Success (200) sample:
```json
{
  "success": true,
  "message": "tenant berhasil diupdate",
  "data": {
    "id": "uuid",
    "name": "Nama Baru",
    "slug": "nama-baru",
    "status": "free",
    "address": "Jl. Sudirman No. 1, Jakarta",
    "description": "Deskripsi kelompok qurban",
    "logo": "https://example.com/logo-baru.png",
    "expired_at": null,
    "owner_id": "uuid",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### 5) Hapus Tenant
- Method: DELETE
- Path: /tenants/:id
- Auth: Bearer token + harus admin tenant
- Description: Soft delete — data tidak langsung dihapus dari DB.
- Success (200): pesan sukses

### 6) Undang Member
- Method: POST
- Path: /tenants/:id/members
- Auth: Bearer token + harus admin tenant
- Description: Mengundang user yang sudah terdaftar di sistem menjadi member dengan role `guest`.
- Request JSON:
```json
{ "user_id": "uuid-user-yang-diundang" }
```
- Success (201) sample:
```json
{
  "success": true,
  "message": "member berhasil diundang",
  "data": {
    "id": "uuid",
    "tenant_id": "uuid",
    "user_id": "uuid",
    "role": "guest",
    "invited_by": "uuid-admin",
    "joined_at": "...",
    "created_at": "...",
    "updated_at": "..."
  }
}
```
- Error 404: user tidak ditemukan
- Error 409: user sudah menjadi member

### 7) List Member Tenant
- Method: GET
- Path: /tenants/:id/members
- Auth: Bearer token + harus member tenant
- Success (200): array member

### 8) Hapus Member
- Method: DELETE
- Path: /tenants/:id/members/:user_id
- Auth: Bearer token + harus admin tenant
- Description: Menghapus member dari tenant. Owner tidak bisa dihapus.
- Success (200): pesan sukses
- Error 403: mencoba hapus owner tenant

---

## Role & Akses

| Role | Hak Akses |
|---|---|
| `admin` | Semua operasi: update, delete, invite, remove member |
| `guest` | Read only: lihat detail dan list member |

- Pembuat tenant otomatis menjadi member dengan role `admin`.
- Owner tenant tidak bisa dihapus dari member list meski oleh admin sekalipun.

---

## Postman Tips
1. Set `base_url` = `http://localhost:8080/api/v1` dan `access_token` dari hasil login.
2. Buat tenant: POST `{{base_url}}/tenants`, simpan `id` dari response ke env var `tenant_id`.
3. Untuk endpoint tenant detail, gunakan `{{base_url}}/tenants/{{tenant_id}}`.
4. Untuk invite member, pastikan `user_id` yang dikirim sudah terdaftar di sistem.
