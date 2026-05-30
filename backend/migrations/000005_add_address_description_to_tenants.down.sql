-- Rollback: Hapus kolom address dan description dari tabel tenants

ALTER TABLE tenants DROP COLUMN IF EXISTS address;
ALTER TABLE tenants DROP COLUMN IF EXISTS description;
