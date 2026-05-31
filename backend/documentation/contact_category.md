# Contact Category API Documentation

Base URL: http://localhost:8080/api/v1

Notes:
- Semua endpoint membutuhkan `Authorization: Bearer <ACCESS_TOKEN>`.
- User harus menjadi member tenant untuk mengakses endpoint ini.
- Operasi create, update, dan delete hanya bisa dilakukan oleh admin tenant.
- Saat tenant baru dibuat, kategori default `umum` dibuat secara otomatis.

---

## Endpoints

### 1) List Kategori Kontak
- Method: GET
- Path: /tenants/:id/contact-categories
- Auth: Bearer token + harus member tenant
- Description: Mengambil semua kategori kontak milik tenant.
- Success (200) sample:
```json
{
  "success": true,
  "message": "daftar kategori kontak berhasil diambil",
  "data": [
    {
      "id": "uuid",
      "tenant_id": "uuid",
      "name": "umum",
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

### 2) Buat Kategori Kontak
- Method: POST
- Path: /tenants/:id/contact-categories
- Auth: Bearer token + harus admin tenant
- Request JSON:
```json
{ "name": "Donatur" }
```
- Success (201) sample:
```json
{
  "success": true,
  "message": "kategori kontak berhasil dibuat",
  "data": {
    "id": "uuid",
    "tenant_id": "uuid",
    "name": "Donatur",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### 3) Update Kategori Kontak
- Method: PUT
- Path: /tenants/:id/contact-categories/:cat_id
- Auth: Bearer token + harus admin tenant
- Request JSON:
```json
{ "name": "Nama Baru" }
```
- Success (200) sample:
```json
{
  "success": true,
  "message": "kategori kontak berhasil diupdate",
  "data": {
    "id": "uuid",
    "tenant_id": "uuid",
    "name": "Nama Baru",
    "created_at": "...",
    "updated_at": "..."
  }
}
```
- Error 404: kategori tidak ditemukan atau bukan milik tenant ini

### 4) Hapus Kategori Kontak
- Method: DELETE
- Path: /tenants/:id/contact-categories/:cat_id
- Auth: Bearer token + harus admin tenant
- Description: Soft delete — data tidak langsung dihapus dari DB.
- Success (200) sample:
```json
{
  "success": true,
  "message": "kategori kontak berhasil dihapus",
  "data": null
}
```
- Error 404: kategori tidak ditemukan atau bukan milik tenant ini

---

## Role & Akses

| Endpoint | Admin | Guest |
|---|---|---|
| GET /contact-categories | ✅ | ✅ |
| POST /contact-categories | ✅ | ❌ |
| PUT /contact-categories/:cat_id | ✅ | ❌ |
| DELETE /contact-categories/:cat_id | ✅ | ❌ |

---

## Postman Tips
1. Set `base_url` = `http://localhost:8080/api/v1` dan `access_token` dari hasil login.
2. Simpan `tenant_id` dari response buat tenant ke env var.
3. Gunakan `{{base_url}}/tenants/{{tenant_id}}/contact-categories`.
4. Simpan `cat_id` dari response buat kategori untuk operasi update/delete.
