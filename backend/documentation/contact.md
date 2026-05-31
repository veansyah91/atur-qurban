# Contact API Documentation

Base URL: http://localhost:8080/api/v1

Notes:
- Semua endpoint membutuhkan `Authorization: Bearer <ACCESS_TOKEN>`.
- User harus menjadi member tenant untuk mengakses endpoint ini.
- Operasi create, update, dan delete hanya bisa dilakukan oleh admin tenant.
- Field `address` dan `phone` bersifat opsional — boleh tidak dikirim atau di-set ke `null`.

---

## Endpoints

### 1) List Kontak
- Method: GET
- Path: /tenants/:id/contacts
- Auth: Bearer token + harus member tenant
- Description: Mengambil semua kontak milik tenant.
- Success (200) sample:
```json
{
  "success": true,
  "message": "daftar kontak berhasil diambil",
  "data": [
    {
      "id": "uuid",
      "tenant_id": "uuid",
      "contact_category_id": "uuid",
      "name": "Ahmad Syukri",
      "is_active": true,
      "address": "Jl. Sudirman No. 1, Jakarta",
      "phone": "+6281234567890",
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

### 2) Buat Kontak
- Method: POST
- Path: /tenants/:id/contacts
- Auth: Bearer token + harus admin tenant
- Request JSON:
```json
{
  "contact_category_id": "uuid",
  "name": "Ahmad Syukri",
  "address": "Jl. Sudirman No. 1, Jakarta",
  "phone": "+6281234567890"
}
```
- Catatan: `address` dan `phone` bersifat **opsional**.
- Success (201) sample:
```json
{
  "success": true,
  "message": "kontak berhasil dibuat",
  "data": {
    "id": "uuid",
    "tenant_id": "uuid",
    "contact_category_id": "uuid",
    "name": "Ahmad Syukri",
    "is_active": true,
    "address": "Jl. Sudirman No. 1, Jakarta",
    "phone": "+6281234567890",
    "created_at": "...",
    "updated_at": "..."
  }
}
```
- Error 404: kategori kontak tidak ditemukan atau bukan milik tenant ini

### 3) Detail Kontak
- Method: GET
- Path: /tenants/:id/contacts/:contact_id
- Auth: Bearer token + harus member tenant
- Success (200): data kontak
- Error 404: kontak tidak ditemukan

### 4) Update Kontak
- Method: PUT
- Path: /tenants/:id/contacts/:contact_id
- Auth: Bearer token + harus admin tenant
- Request JSON:
```json
{
  "contact_category_id": "uuid",
  "name": "Nama Baru",
  "is_active": true,
  "address": "Jl. Merdeka No. 5, Bandung",
  "phone": "+6289876543210"
}
```
- Catatan: `is_active` harus selalu dikirim (true/false). `address` dan `phone` opsional.
- Success (200) sample:
```json
{
  "success": true,
  "message": "kontak berhasil diupdate",
  "data": {
    "id": "uuid",
    "tenant_id": "uuid",
    "contact_category_id": "uuid",
    "name": "Nama Baru",
    "is_active": true,
    "address": "Jl. Merdeka No. 5, Bandung",
    "phone": "+6289876543210",
    "created_at": "...",
    "updated_at": "..."
  }
}
```
- Error 404: kontak atau kategori tidak ditemukan

### 5) Hapus Kontak
- Method: DELETE
- Path: /tenants/:id/contacts/:contact_id
- Auth: Bearer token + harus admin tenant
- Description: Soft delete — data tidak langsung dihapus dari DB.
- Success (200) sample:
```json
{
  "success": true,
  "message": "kontak berhasil dihapus",
  "data": null
}
```
- Error 404: kontak tidak ditemukan

---

## Role & Akses

| Endpoint | Admin | Guest |
|---|---|---|
| GET /contacts | ✅ | ✅ |
| POST /contacts | ✅ | ❌ |
| GET /contacts/:contact_id | ✅ | ✅ |
| PUT /contacts/:contact_id | ✅ | ❌ |
| DELETE /contacts/:contact_id | ✅ | ❌ |

---

## Postman Tips
1. Set `base_url` = `http://localhost:8080/api/v1` dan `access_token` dari hasil login.
2. Pastikan `tenant_id` dan `cat_id` sudah tersedia di env var.
3. Gunakan `{{base_url}}/tenants/{{tenant_id}}/contacts`.
4. Untuk operasi pada kontak tertentu, tambahkan `/:contact_id` di akhir path.
