-- ============================================================================
-- 009 — Tabel: payments
-- ============================================================================
-- APA ITU?
--   Tabel payments menyimpan data PEMBAYARAN untuk setiap order.
--   Setiap order bisa punya beberapa percobaan pembayaran (retry).
--
-- FUNGSI:
--   - Menyimpan metode pembayaran (cash/QRIS/transfer)
--   - Tracking status pembayaran: pending → paid / failed / refunded
--   - Idempotency key mencegah double payment
--   - Reference number untuk tracking dari payment gateway
--   - Data ini dipakai untuk laporan revenue merchant
--
-- CONTOH DATA:
--   | order_id | tenant_id | method | status | amount   | paid_at              |
--   |----------|-----------|--------|--------|----------|----------------------|
--   | uuid-ord | uuid-ten  | cash   | paid   | 75000.00 | 2026-04-18 12:30:00  |
--   | uuid-ord | uuid-ten  | qris   | failed | 75000.00 | NULL                 |
--
-- RELASI:
--   payments (N) → (1) orders   — setiap payment untuk satu order
--                                  ON DELETE RESTRICT: jangan hapus order yang ada payment
--   payments (N) → (1) tenants  — setiap payment milik satu tenant
--                                  ON DELETE RESTRICT: jangan hapus tenant dengan transaksi
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS payments (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke orders: payment ini untuk order mana
    -- ON DELETE RESTRICT: JANGAN hapus order yang sudah ada payment-nya
    -- (data keuangan harus selalu dijaga)
    order_id        UUID           NOT NULL
                    REFERENCES orders(id) ON DELETE RESTRICT,

    -- FK ke tenants: payment ini milik tenant mana
    -- ON DELETE RESTRICT: JANGAN hapus tenant yang punya data keuangan
    tenant_id       UUID           NOT NULL
                    REFERENCES tenants(id) ON DELETE RESTRICT,

    -- Idempotency key: mencegah double payment karena network retry
    -- Client generate UUID → kirim bersama request → jika sudah ada → 409 Conflict
    idempotency_key VARCHAR(120)   UNIQUE,

    -- Metode pembayaran:
    --   cash     = bayar tunai di kasir
    --   qris     = scan QRIS (e-wallet / m-banking)
    --   transfer = transfer bank
    method          VARCHAR(20)    NOT NULL
                    CHECK (method IN ('cash', 'qris', 'transfer')),

    -- Status lifecycle payment:
    --   pending  → pembayaran dibuat, menunggu konfirmasi
    --   paid     → pembayaran berhasil / terkonfirmasi
    --   failed   → pembayaran gagal
    --   refunded → pembayaran sudah di-refund ke customer
    status          VARCHAR(20)    NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'paid', 'failed', 'refunded')),

    -- Jumlah yang harus dibayar (harus sama dengan orders.total)
    amount          NUMERIC(12,2)  NOT NULL CHECK (amount >= 0),

    -- Timestamp kapan pembayaran berhasil (null jika belum paid)
    paid_at         TIMESTAMPTZ,

    -- Nomor referensi dari payment gateway (opsional)
    -- Contoh: "TXN-QRIS-20260418-001" dari provider QRIS
    reference_no    VARCHAR(120),

    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ             -- soft delete
);

-- Index: cari payment berdasarkan order
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments (order_id) WHERE deleted_at IS NULL;

-- Index: filter payment per tenant
CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments (tenant_id) WHERE deleted_at IS NULL;

-- Index: sort payment terbaru per tenant (untuk laporan keuangan)
CREATE INDEX IF NOT EXISTS idx_payments_tenant_created_at ON payments (tenant_id, created_at DESC) WHERE deleted_at IS NULL;

-- Index: filter payment berdasarkan status (untuk monitoring pending payments)
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
