-- Tujuan: Menyimpan data customer yang terkait dengan order pada tenant tertentu.

BEGIN;

CREATE TABLE IF NOT EXISTS customers (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID         NOT NULL,
    tenant_id     UUID         NOT NULL,
    full_name     VARCHAR(150) NOT NULL,
    phone_number  VARCHAR(20),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT fk_customers_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_customers_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_customers_order_id ON customers (order_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers (tenant_id) WHERE deleted_at IS NULL;

COMMIT;