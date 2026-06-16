BEGIN;

-- =====================================================
-- 1. Tambah kolom payment_channel (Nullable by default)
-- =====================================================
ALTER TABLE orders
ADD COLUMN payment_channel VARCHAR(50);

-- =====================================================
-- 2. Migrasi data existing berdasarkan payment_method
--    - CASH -> 'OFFLINE'
--    - Selain CASH -> default awal sesuai payment_method
-- =====================================================
UPDATE orders
SET payment_channel = CASE 
    WHEN payment_method = 'CASH' THEN 'OFFLINE'
    WHEN payment_method = 'TRANSFER_BANK' THEN 'BANK_TRANSFER'
    WHEN payment_method = 'QRIS' THEN 'QRIS'
    WHEN payment_method = 'E_WALLET' THEN 'E_WALLET'
    WHEN payment_method = 'KARTU_KREDIT' THEN 'CREDIT_CARD'
    WHEN payment_method = 'MINIMARKET' THEN 'MINIMARKET'
    ELSE NULL
END;

-- =====================================================
-- 3. Tambah validasi (Constraint)
-- =====================================================
ALTER TABLE orders
ADD CONSTRAINT chk_orders_payment_channel
CHECK (
    (payment_method IS NULL AND payment_channel IS NULL) OR
    (payment_method = 'CASH' AND payment_channel = 'OFFLINE') OR
    (payment_method = 'QRIS' AND (payment_channel IS NULL OR payment_channel = 'QRIS')) OR
    (payment_method = 'E_WALLET' AND (payment_channel IS NULL OR payment_channel IN ('GOPAY', 'SHOPEEPAY', 'SHOPEPAY', 'DANA', 'OVO', 'E_WALLET'))) OR
    (payment_method = 'TRANSFER_BANK' AND (payment_channel IS NULL OR payment_channel IN ('BCA', 'BRI', 'BNI', 'MANDIRI', 'PERMATA', 'DANAMON', 'BSI', 'SEABANK', 'CIMB', 'BANK_TRANSFER'))) OR
    (payment_method = 'KARTU_KREDIT' AND (payment_channel IS NULL OR payment_channel IN ('AKULAKU', 'AKULAKUPAYLATER', 'KREDIVO', 'CREDIT_CARD'))) OR
    (payment_method = 'MINIMARKET' AND (payment_channel IS NULL OR payment_channel IN ('INDOMARET', 'ALFAMART', 'ALFAMIDI', 'MINIMARKET')))
);

COMMIT;
