-- Tujuan: Menghapus tabel customers karena data customer sudah didelete/tidak digunakan lagi.
BEGIN;

DROP TABLE IF EXISTS customers CASCADE;

COMMIT;
