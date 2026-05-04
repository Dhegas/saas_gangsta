-- Tujuan: Pindahkan data tenant_profiles ke tenants dan lengkapi kolom baru.

BEGIN;

ALTER TABLE tenants
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20),
    ADD COLUMN IF NOT EXISTS open_hours VARCHAR(100),
    ADD COLUMN IF NOT EXISTS logo_url TEXT,
    ADD COLUMN IF NOT EXISTS banner_url TEXT;

WITH ranked_profiles AS (
    SELECT
        tp.tenant_id,
        tp.name,
        tp.description,
        tp.updated_at,
        ROW_NUMBER() OVER (
            PARTITION BY tp.tenant_id
            ORDER BY tp.is_active DESC, tp.sort_order ASC, tp.updated_at DESC, tp.created_at DESC
        ) AS rn
    FROM tenant_profiles tp
    WHERE tp.deleted_at IS NULL
)
UPDATE tenants t
SET name = rp.name,
    description = rp.description,
    updated_at = GREATEST(t.updated_at, rp.updated_at)
FROM ranked_profiles rp
WHERE rp.rn = 1
  AND rp.tenant_id = t.id;

COMMIT;
