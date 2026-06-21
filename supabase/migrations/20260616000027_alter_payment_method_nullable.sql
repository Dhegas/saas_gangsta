BEGIN;

-- =====================================================
-- 1. Hapus NOT NULL dari kolom payment_method
-- =====================================================
ALTER TABLE orders
ALTER COLUMN payment_method DROP NOT NULL;

COMMIT;
