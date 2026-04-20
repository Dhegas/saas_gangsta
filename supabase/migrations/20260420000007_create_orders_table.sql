-- Tujuan: Menyimpan transaksi pesanan customer per tenant.

BEGIN;

CREATE TABLE IF NOT EXISTS orders (
    id                UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID          NOT NULL,
    user_id           UUID,
    dining_tables_id  UUID,
    status            VARCHAR(20)   NOT NULL,
    total_price       NUMERIC(12,2) NOT NULL,
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT fk_orders_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL,
    CONSTRAINT fk_orders_dining_table
        FOREIGN KEY (dining_tables_id)
        REFERENCES dining_tables(id)
        ON DELETE SET NULL,
    CONSTRAINT chk_orders_total_price_non_negative CHECK (total_price >= 0)
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_dining_tables_id ON orders (dining_tables_id) WHERE deleted_at IS NULL;

-- Index kolom yang sering difilter.
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status) WHERE deleted_at IS NULL;

COMMIT;