-- Tujuan: Menyimpan daftar meja makan per tenant.

BEGIN;

CREATE TABLE IF NOT EXISTS dining_tables (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID         NOT NULL,
    table_name  VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_dining_tables_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_dining_tables_tenant_id ON dining_tables (tenant_id) WHERE deleted_at IS NULL;

COMMIT;