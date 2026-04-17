-- ============================================================================
-- 001 — Tabel: tenants
-- ============================================================================
-- APA ITU?
--   Tabel tenants adalah ROOT ENTITY (entitas utama) dalam sistem multi-tenant.
--   Setiap "tenant" merepresentasikan SATU TOKO/MERCHANT yang terdaftar
--   di platform SaaS ini.
--
-- FUNGSI:
--   - Menyimpan data identitas toko (nama, slug URL, status aktif/nonaktif)
--   - Menjadi parent (induk) dari SEMUA tabel bisnis lainnya
--   - Setiap data bisnis (menu, order, payment, dll) WAJIB terhubung ke tenant
--     melalui kolom tenant_id — ini yang disebut "tenant isolation"
--
-- CONTOH DATA:
--   | id   | name           | slug          | status   |
--   |------|----------------|---------------|----------|
--   | uuid | Warung Makan A | warung-makan-a| active   |
--   | uuid | Kafe Kopi B    | kafe-kopi-b   | suspended|
--
-- RELASI:
--   tenants (1) → (N) users
--   tenants (1) → (1) merchant_profiles
--   tenants (1) → (N) categories
--   tenants (1) → (N) menus
--   tenants (1) → (N) tables
--   tenants (1) → (N) orders
--   tenants (1) → (N) payments
--   tenants (1) → (N) subscriptions
--   tenants (1) → (N) reviews
--   tenants (1) → (N) audit_logs
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS tenants (
    -- Primary key: UUID di-generate otomatis oleh database
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Nama toko/brand yang ditampilkan di UI
    name        VARCHAR(150) NOT NULL,

    -- Slug: versi URL-friendly dari nama toko (unik di seluruh platform)
    -- Contoh: "warung-makan-a" → digunakan di URL public
    slug        VARCHAR(80)  NOT NULL UNIQUE,

    -- Status tenant:
    --   active    = toko beroperasi normal
    --   inactive  = toko non-aktif (belum bayar / tutup sementara)
    --   suspended = di-suspend oleh admin platform (pelanggaran, dll)
    status      VARCHAR(20)  NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'inactive', 'suspended')),

    -- Timestamp kapan data dibuat
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- Timestamp kapan data terakhir diubah (auto-update via trigger)
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- Soft delete: jika diisi, berarti tenant sudah "dihapus" tapi datanya
    -- masih ada di database (tidak benar-benar dihapus)
    deleted_at  TIMESTAMPTZ
);

-- Index untuk pencarian berdasarkan slug (sering dipakai di lookup URL)
CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants (slug);

-- Index untuk filter status, hanya untuk tenant yang belum di-soft-delete
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants (status) WHERE deleted_at IS NULL;

-- Trigger: otomatis update kolom updated_at setiap kali baris di-UPDATE
CREATE TRIGGER trg_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
