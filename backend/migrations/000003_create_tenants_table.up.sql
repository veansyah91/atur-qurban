-- Migration: Create tenants table
-- Tabel untuk menyimpan data tenant/organisasi

CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'free' CHECK (status IN ('free', 'paid')),
    expired_at TIMESTAMP,
    owner_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index untuk soft delete queries
CREATE INDEX IF NOT EXISTS idx_tenants_deleted_at ON tenants(deleted_at);

-- Index untuk owner lookup
CREATE INDEX IF NOT EXISTS idx_tenants_owner_id ON tenants(owner_id);

-- Index untuk slug lookup
CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants(slug);

-- Index untuk status lookup
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);
