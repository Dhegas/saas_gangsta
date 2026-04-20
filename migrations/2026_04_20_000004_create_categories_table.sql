-- Tujuan: Menyimpan kategori menu per tenant agar menu terstruktur.

BEGIN;

CREATE TABLE IF NOT EXISTS categories (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID         NOT NULL,
    name        VARCHAR(120) NOT NULL,
    description TEXT,
    sort_order  INTEGER      NOT NULL DEFAULT 0,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT uq_categories_tenant_name UNIQUE (tenant_id, name),
    CONSTRAINT fk_categories_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories (tenant_id) WHERE deleted_at IS NULL;

-- Index kolom yang sering difilter.
CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories (is_active) WHERE deleted_at IS NULL;

COMMIT;