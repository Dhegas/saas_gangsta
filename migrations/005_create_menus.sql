-- ============================================================================
-- 005 — Tabel: menus
-- ============================================================================
-- APA ITU?
--   Tabel menus menyimpan DAFTAR MENU (makanan/minuman) per toko/tenant.
--   Ini adalah data utama yang dilihat customer saat scan QR code.
--
-- FUNGSI:
--   - Menyimpan informasi menu: nama, harga, deskripsi, gambar
--   - Merchant bisa kelola (CRUD) dan toggle ketersediaan (habis/tersedia)
--   - Customer melihat menu yang is_available = true
--   - Harga dan nama menu di-snapshot ke order_items saat order dibuat
--     (agar riwayat order tetap akurat meski menu diubah di kemudian hari)
--
-- CONTOH DATA:
--   | tenant_id | category_id | name           | price    | is_available |
--   |-----------|-------------|----------------|----------|--------------|
--   | uuid-123  | uuid-cat-1  | Nasi Goreng    | 25000.00 | true         |
--   | uuid-123  | uuid-cat-2  | Es Teh Manis   | 8000.00  | true         |
--   | uuid-123  | uuid-cat-1  | Mie Ayam       | 20000.00 | false        |
--
-- RELASI:
--   menus (N) → (1) tenants      — setiap menu milik satu tenant
--   menus (N) → (1) categories   — setiap menu bisa masuk satu kategori (opsional)
--   menus (1) → (N) order_items  — satu menu bisa dipesan di banyak order
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS menus (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: menu ini milik tenant mana
    -- ON DELETE CASCADE: jika tenant dihapus, semua menu-nya juga dihapus
    tenant_id       UUID           NOT NULL
                    REFERENCES tenants(id) ON DELETE CASCADE,

    -- FK ke categories: menu ini masuk kategori mana
    -- Nullable: menu boleh tidak punya kategori
    -- ON DELETE SET NULL: jika kategori dihapus, menu tetap ada tapi category_id jadi NULL
    category_id     UUID
                    REFERENCES categories(id) ON DELETE SET NULL,

    -- Nama menu yang ditampilkan ke customer
    name            VARCHAR(180)   NOT NULL,

    -- Deskripsi menu (bahan, porsi, dll)
    description     TEXT,

    -- Harga per item (dalam Rupiah, precision 12 digit, 2 desimal)
    -- CHECK: harga tidak boleh negatif
    price           NUMERIC(12,2)  NOT NULL CHECK (price >= 0),

    -- URL gambar menu (disimpan di storage)
    image_url       TEXT,

    -- Toggle tersedia/habis
    -- true  = tampil di menu customer dan bisa dipesan
    -- false = tidak tampil / habis
    is_available    BOOLEAN        NOT NULL DEFAULT TRUE,

    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ             -- soft delete
);

-- Index: filter menu berdasarkan tenant
CREATE INDEX IF NOT EXISTS idx_menus_tenant_id ON menus (tenant_id) WHERE deleted_at IS NULL;

-- Index: filter menu berdasarkan kategori (untuk tampilan per kategori)
CREATE INDEX IF NOT EXISTS idx_menus_category_id ON menus (category_id) WHERE deleted_at IS NULL;

-- Index: filter menu yang tersedia per tenant (query paling sering dari customer)
CREATE INDEX IF NOT EXISTS idx_menus_tenant_available ON menus (tenant_id, is_available) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_menus_updated_at
    BEFORE UPDATE ON menus
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
