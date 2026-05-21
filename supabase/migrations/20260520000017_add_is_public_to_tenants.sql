-- Tujuan: Menambahkan kolom is_public pada tabel tenants.

BEGIN;

ALTER TABLE tenants
    ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT true;

-- Update data yang sudah ada agar nilainya default true
UPDATE tenants
SET is_public = true
WHERE is_public IS NULL;

-- Index untuk filter query tenant public
CREATE INDEX IF NOT EXISTS idx_tenants_is_public ON tenants (is_public) WHERE deleted_at IS NULL;

COMMIT;
