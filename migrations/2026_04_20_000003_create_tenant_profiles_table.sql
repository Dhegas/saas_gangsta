-- Tujuan: Menyimpan profil tenant (restoran/UMKM) untuk kebutuhan branding dan informasi usaha.

BEGIN;

CREATE TABLE IF NOT EXISTS tenant_profiles (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID         NOT NULL,
    name        VARCHAR(120) NOT NULL UNIQUE,
    description TEXT,
    sort_order  INTEGER      NOT NULL DEFAULT 0,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_tenant_profiles_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_tenant_profiles_tenant_id ON tenant_profiles (tenant_id) WHERE deleted_at IS NULL;

-- Index kolom yang sering difilter.
CREATE INDEX IF NOT EXISTS idx_tenant_profiles_is_active ON tenant_profiles (is_active) WHERE deleted_at IS NULL;

COMMIT;