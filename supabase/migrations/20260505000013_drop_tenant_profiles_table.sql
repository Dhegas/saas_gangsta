-- Tujuan: Hapus tabel tenant_profiles setelah data dipindahkan.

BEGIN;

DROP TABLE IF EXISTS tenant_profiles;

COMMIT;
