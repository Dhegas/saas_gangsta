BEGIN;

-- =====================================================
-- 1. Tambah kolom baru
-- =====================================================

ALTER TABLE orders
ADD COLUMN queue_number VARCHAR(50),
ADD COLUMN payment_method VARCHAR(50);

-- =====================================================
-- 2. Isi data existing
-- =====================================================

UPDATE orders
SET payment_method = 'CASH'
WHERE payment_method IS NULL;

-- =====================================================
-- 3. Wajibkan payment_method
-- =====================================================

ALTER TABLE orders
ALTER COLUMN payment_method SET NOT NULL;

-- =====================================================
-- 4. Validasi metode pembayaran
-- =====================================================

ALTER TABLE orders
ADD CONSTRAINT chk_orders_payment_method
CHECK (
    payment_method IN (
        'QRIS',
        'TRANSFER_BANK',
        'CASH',
        'E_WALLET',
        'KARTU_KREDIT',
        'MINIMARKET'
    )
);

-- =====================================================
-- 5. Cegah nomor antrian duplikat
--    dalam tenant dan hari yang sama
-- =====================================================

CREATE UNIQUE INDEX idx_unique_tenant_queue_daily
ON orders (
    tenant_id,
    queue_number,
    (timezone('UTC', created_at)::date)
)
WHERE deleted_at IS NULL
AND queue_number IS NOT NULL;

COMMIT;