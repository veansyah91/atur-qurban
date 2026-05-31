CREATE TABLE contact_categories (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_contact_categories_tenant_id ON contact_categories(tenant_id);
CREATE INDEX idx_contact_categories_deleted_at ON contact_categories(deleted_at);
