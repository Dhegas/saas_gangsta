-- ============================================================================
-- 008 — Tabel: order_items
-- ============================================================================
-- APA ITU?
--   Tabel order_items menyimpan DETAIL ITEM di dalam setiap order.
--   Ini adalah tabel penghubung antara orders dan menus, yang menyimpan
--   jumlah pesanan, harga satuan, dan subtotal per item.
--
-- FUNGSI:
--   - Menyimpan item apa saja yang dipesan dalam satu order
--   - Menyimpan SNAPSHOT harga dan nama menu SAAT ORDER DIBUAT
--     (penting! karena harga/nama menu bisa berubah di kemudian hari,
--      tapi riwayat order harus tetap menunjukkan harga asli saat pesan)
--   - Mendukung catatan per item (contoh: "tanpa sambal", "extra keju")
--
-- CONTOH DATA:
--   | order_id | menu_id  | menu_name     | qty | unit_price | subtotal  | notes       |
--   |----------|----------|---------------|-----|------------|-----------|-------------|
--   | uuid-ord | uuid-mn1 | Nasi Goreng   | 2   | 25000.00   | 50000.00  | extra telur |
--   | uuid-ord | uuid-mn2 | Es Teh Manis  | 3   | 8000.00    | 24000.00  | less sugar  |
--
-- RELASI:
--   order_items (N) → (1) orders — setiap item milik satu order
--                                   ON DELETE CASCADE: hapus order = hapus semua itemnya
--   order_items (N) → (1) menus  — setiap item mereferensi satu menu
--                                   ON DELETE RESTRICT: jangan hapus menu yang masih direferensi
--
-- CATATAN:
--   Tabel ini TIDAK punya updated_at dan deleted_at karena:
--   - Item order bersifat immutable (tidak boleh diubah setelah dibuat)
--   - Jika order dihapus/cancel, semua item ikut dihapus via CASCADE
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS order_items (
    id          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke orders: item ini milik order mana
    -- ON DELETE CASCADE: jika order dihapus, semua item-nya otomatis ikut dihapus
    order_id    UUID           NOT NULL
                REFERENCES orders(id) ON DELETE CASCADE,

    -- FK ke menus: item ini mereferensi menu mana
    -- ON DELETE RESTRICT: JANGAN hapus menu yang masih ada di order aktif
    menu_id     UUID           NOT NULL
                REFERENCES menus(id) ON DELETE RESTRICT,

    -- Snapshot nama menu saat order dibuat
    -- Kenapa disimpan? Karena merchant bisa ubah nama menu kapan saja,
    -- tapi riwayat order harus tetap menunjukkan nama asli saat order
    menu_name   VARCHAR(180),

    -- Jumlah item yang dipesan (minimal 1)
    quantity    INTEGER        NOT NULL CHECK (quantity > 0),

    -- Snapshot harga satuan SAAT ORDER DIBUAT (bukan harga menu saat ini)
    -- Ini memastikan laporan revenue tetap akurat
    unit_price  NUMERIC(12,2)  NOT NULL CHECK (unit_price >= 0),

    -- Subtotal = quantity × unit_price (dihitung oleh backend)
    subtotal    NUMERIC(12,2)  NOT NULL CHECK (subtotal >= 0),

    -- Catatan khusus per item (contoh: "tanpa sambal", "pedas level 5")
    notes       TEXT
);

-- Index: cari semua item dalam satu order (JOIN query paling sering)
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);

-- Index: cari semua order yang berisi menu tertentu (untuk laporan populer menu)
CREATE INDEX IF NOT EXISTS idx_order_items_menu_id ON order_items (menu_id);

COMMIT;
