-- Migration: Add verified_at column to users table
-- Gunakan untuk menandai kapan user melakukan verifikasi nomor telepon

ALTER TABLE users ADD COLUMN verified_at TIMESTAMP DEFAULT NULL;

-- Index untuk query user yang sudah terverifikasi
CREATE INDEX IF NOT EXISTS idx_users_verified_at ON users(verified_at);
