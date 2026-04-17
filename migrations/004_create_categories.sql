-- ============================================================================
-- 004 — Tabel: categories
-- ============================================================================
-- APA ITU?
--   Tabel categories menyimpan KATEGORI MENU per toko/tenant.
--   Contoh kategori: "Makanan Berat", "Minuman", "Snack", "Dessert".
--
-- FUNGSI:
--   - Mengelompokkan menu items agar tampilan menu lebih rapi
--   - Customer bisa filter menu berdasarkan kategori
--   - Merchant bisa atur urutan tampil (sort_order) dan aktif/nonaktif
--
-- CONTOH DATA:
--   | tenant_id | name          | sort_order | is_active |
--   |-----------|---------------|------------|-----------|
--   | uuid-123  | Makanan Berat | 1          | true      |
--   | uuid-123  | Minuman       | 2          | true      |
--   | uuid-123  | Promo Spesial | 3          | false     |
--
-- RELASI:
--   categories (N) → (1) tenants  — setiap kategori milik satu tenant
--   categories (1) → (N) menus    — satu kategori bisa punya banyak menu
--
-- CONSTRAINT UNIK:
--   (tenant_id, name) — nama kategori tidak boleh duplikat dalam satu tenant
--   Contoh: Tenant A boleh punya "Minuman", Tenant B juga boleh punya "Minuman"
--           tapi Tenant A TIDAK boleh punya dua kategori bernama "Minuman"
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS categories (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: kategori ini milik tenant mana
    -- ON DELETE CASCADE: jika tenant dihapus, semua kategorinya juga dihapus
    tenant_id   UUID         NOT NULL
                REFERENCES tenants(id) ON DELETE CASCADE,

    -- Nama kategori (unik per tenant)
    name        VARCHAR(120) NOT NULL,

    -- Deskripsi opsional untuk kategori
    description TEXT,

    -- Urutan tampil di UI (ascending: 1 tampil pertama, 2 kedua, dst)
    sort_order  INTEGER      NOT NULL DEFAULT 0,

    -- Toggle aktif/nonaktif kategori (jika false, tidak tampil di menu customer)
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,          -- soft delete

    -- Constraint: nama kategori unik dalam satu tenant
    UNIQUE (tenant_id, name)
);

-- Index: filter kategori berdasarkan tenant (dipakai di hampir semua query)
CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories (tenant_id) WHERE deleted_at IS NULL;

-- Index: filter kategori aktif per tenant (dipakai saat customer lihat menu)
CREATE INDEX IF NOT EXISTS idx_categories_tenant_active ON categories (tenant_id, is_active) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
