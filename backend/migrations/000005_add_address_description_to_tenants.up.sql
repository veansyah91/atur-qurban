-- Migration: Tambah kolom address dan description ke tabel tenants

ALTER TABLE tenants ADD COLUMN IF NOT EXISTS address VARCHAR(500);
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS description TEXT;
