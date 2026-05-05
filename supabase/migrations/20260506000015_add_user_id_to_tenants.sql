-- Tujuan: Pindahkan relasi user-tenant ke tabel tenants lewat kolom user_id.
-- Catatan: Menghapus users.tenant_id setelah data dipindahkan.

BEGIN;

ALTER TABLE tenants
    ADD COLUMN IF NOT EXISTS user_id UUID;

ALTER TABLE tenants
    ADD CONSTRAINT fk_tenants_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL;

-- Index untuk query tenant berdasarkan owner/user.
CREATE INDEX IF NOT EXISTS idx_tenants_user_id ON tenants (user_id) WHERE deleted_at IS NULL;

-- Migrasi data lama: user.tenant_id -> tenants.user_id (ambil user paling awal jika ada duplikasi).
WITH owner_map AS (
    SELECT DISTINCT ON (tenant_id)
        tenant_id,
        id AS user_id
    FROM users
    WHERE tenant_id IS NOT NULL
    ORDER BY tenant_id, created_at ASC
)
UPDATE tenants t
SET user_id = om.user_id
FROM owner_map om
WHERE t.id = om.tenant_id;

-- Hapus relasi lama di users.
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_tenant;
DROP INDEX IF EXISTS idx_users_tenant_id;
ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;

COMMIT;
