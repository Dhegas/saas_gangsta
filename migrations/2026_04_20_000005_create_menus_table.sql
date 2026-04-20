-- Tujuan: Menyimpan data menu makanan/minuman per tenant.

BEGIN;

CREATE TABLE IF NOT EXISTS menus (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID          NOT NULL,
    category_id   UUID,
    name          VARCHAR(180)  NOT NULL,
    description   TEXT,
    price         NUMERIC(12,2) NOT NULL,
    image_url     TEXT,
    is_available  BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT fk_menus_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_menus_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE SET NULL,
    CONSTRAINT chk_menus_price_non_negative CHECK (price >= 0)
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_menus_tenant_id ON menus (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_menus_category_id ON menus (category_id) WHERE deleted_at IS NULL;

-- Index kolom yang sering difilter.
CREATE INDEX IF NOT EXISTS idx_menus_is_available ON menus (is_available) WHERE deleted_at IS NULL;

COMMIT;