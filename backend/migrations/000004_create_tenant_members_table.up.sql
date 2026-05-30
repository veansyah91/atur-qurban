-- Migration: Create tenant_members table
-- Tabel untuk menyimpan data keanggotaan user di dalam tenant

CREATE TABLE IF NOT EXISTS tenant_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'guest' CHECK (role IN ('admin', 'guest')),
    invited_by UUID REFERENCES users(id),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(tenant_id, user_id)
);

-- Index untuk soft delete queries
CREATE INDEX IF NOT EXISTS idx_tenant_members_deleted_at ON tenant_members(deleted_at);

-- Index untuk tenant lookup
CREATE INDEX IF NOT EXISTS idx_tenant_members_tenant_id ON tenant_members(tenant_id);

-- Index untuk user lookup
CREATE INDEX IF NOT EXISTS idx_tenant_members_user_id ON tenant_members(user_id);

-- Index untuk role lookup
CREATE INDEX IF NOT EXISTS idx_tenant_members_role ON tenant_members(role);
