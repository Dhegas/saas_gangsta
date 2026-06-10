BEGIN;

-- =====================================================
-- Tambahkan kolom payment ke tabel orders
-- =====================================================

ALTER TABLE orders
    ADD COLUMN payment_status   VARCHAR(20)  NOT NULL DEFAULT 'UNPAID',
    ADD COLUMN midtrans_order_id VARCHAR(255) NULL,
    ADD COLUMN midtrans_transaction_id VARCHAR(255) NULL,
    ADD COLUMN paid_at           TIMESTAMPTZ  NULL;

-- =====================================================
-- Constraint: payment_status hanya boleh nilai enum ini
-- =====================================================

ALTER TABLE orders
    ADD CONSTRAINT chk_orders_payment_status
    CHECK (payment_status IN ('UNPAID', 'PAID', 'FAILED', 'REFUNDED'));

-- =====================================================
-- Unique: satu midtrans_order_id hanya boleh untuk satu order
-- =====================================================

CREATE UNIQUE INDEX idx_orders_midtrans_order_id
    ON orders (midtrans_order_id)
    WHERE midtrans_order_id IS NOT NULL;

-- =====================================================
-- Index untuk lookup cepat saat webhook datang
-- =====================================================

CREATE INDEX idx_orders_payment_status
    ON orders (payment_status)
    WHERE deleted_at IS NULL;

COMMIT;
