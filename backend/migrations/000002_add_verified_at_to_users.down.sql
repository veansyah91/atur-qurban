-- Rollback: Remove verified_at column dari users table

ALTER TABLE users DROP COLUMN IF EXISTS verified_at;

-- Hapus index
DROP INDEX IF EXISTS idx_users_verified_at;
