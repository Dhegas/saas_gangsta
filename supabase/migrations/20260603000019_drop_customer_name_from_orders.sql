-- Tujuan: Menghapus kolom customer_name dari tabel orders karena data customer dibaca dari relasi tabel users.
BEGIN;

ALTER TABLE IF EXISTS orders DROP COLUMN IF EXISTS customer_name;

COMMIT;
