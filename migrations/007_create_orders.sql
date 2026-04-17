-- ============================================================================
-- 007 — Tabel: orders
-- ============================================================================
-- APA ITU?
--   Tabel orders menyimpan data PESANAN dari customer (self-order via QR)
--   maupun dari merchant (manual via POS/kasir).
--
-- FUNGSI:
--   - Menyimpan header order: siapa yang pesan, dari meja mana, total harga
--   - Mendukung lifecycle STATUS ORDER (state machine):
--       pending → accepted → cooking → ready → done
--       pending → canceled (bisa dibatalkan sebelum accepted)
--   - Idempotency key mencegah duplicate order dari double-click/retry
--   - Order number untuk nomor struk yang tampil di receipt
--   - Bisa tracking asal order: dari self_order atau pos
--
-- CONTOH DATA:
--   | tenant_id | table_id | user_id  | status  | total    | order_source |
--   |-----------|----------|----------|---------|----------|--------------|
--   | uuid-123  | uuid-tbl | uuid-usr | pending | 75000.00 | self_order   |
--   | uuid-123  | NULL     | NULL     | cooking | 50000.00 | pos          |
--
-- RELASI:
--   orders (N) → (1) tenants   — setiap order milik satu tenant
--   orders (N) → (1) tables    — order dari meja mana (nullable, POS bisa tanpa meja)
--   orders (N) → (1) users     — customer yang pesan (nullable untuk guest order)
--   orders (1) → (N) order_items — satu order berisi banyak item menu
--   orders (1) → (N) payments   — satu order bisa punya banyak payment attempt
--   orders (1) → (1) reviews    — satu order hanya bisa di-review sekali
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS orders (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: order ini milik tenant mana
    -- ON DELETE RESTRICT: JANGAN hapus tenant yang masih punya order
    -- (data transaksi harus dijaga, jadi pakai RESTRICT bukan CASCADE)
    tenant_id       UUID           NOT NULL
                    REFERENCES tenants(id) ON DELETE RESTRICT,

    -- FK ke tables: dari meja mana order ini
    -- Nullable: order dari POS bisa tanpa meja
    -- ON DELETE SET NULL: jika meja dihapus, order tetap ada tapi table_id = NULL
    table_id        UUID
                    REFERENCES tables(id) ON DELETE SET NULL,

    -- FK ke users: customer yang membuat order
    -- Nullable: guest order (tanpa login) diperbolehkan
    -- ON DELETE SET NULL: jika user dihapus, order tetap ada
    user_id         UUID
                    REFERENCES users(id) ON DELETE SET NULL,

    -- Nomor order/struk unik per tenant (di-generate oleh backend)
    -- Contoh: "ORD-20260418-001"
    order_number    VARCHAR(50),

    -- Idempotency key: client-generated UUID untuk mencegah duplicate order
    -- Jika client kirim request yang sama 2x, yang kedua akan ditolak (409 Conflict)
    idempotency_key VARCHAR(120)   UNIQUE,

    -- Status lifecycle order (state machine):
    --   pending  → order baru masuk, belum dikonfirmasi merchant
    --   accepted → merchant terima order, mulai diproses
    --   cooking  → sedang dimasak di dapur
    --   ready    → makanan siap disajikan
    --   done     → order selesai, sudah diserahkan ke customer
    --   canceled → order dibatalkan
    status          VARCHAR(20)    NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'accepted', 'cooking', 'ready', 'done', 'canceled')),

    -- Subtotal: jumlah harga semua item sebelum pajak
    subtotal        NUMERIC(12,2)  NOT NULL DEFAULT 0,

    -- Pajak (jika ada)
    tax             NUMERIC(12,2)  NOT NULL DEFAULT 0,

    -- Total: subtotal + tax (angka final yang harus dibayar)
    total           NUMERIC(12,2)  NOT NULL DEFAULT 0,

    -- Catatan umum untuk order (contoh: "pesanan meja 3, gak pake nasi")
    notes           TEXT,

    -- Asal order:
    --   self_order = customer pesan sendiri via QR scan
    --   pos        = merchant input manual via kasir/POS
    order_source    VARCHAR(20)    NOT NULL
                    CHECK (order_source IN ('self_order', 'pos')),

    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ             -- soft delete
);

-- Index: filter order berdasarkan tenant
CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders (tenant_id) WHERE deleted_at IS NULL;

-- Index: sort order terbaru per tenant (untuk dashboard merchant)
CREATE INDEX IF NOT EXISTS idx_orders_tenant_created_at ON orders (tenant_id, created_at DESC) WHERE deleted_at IS NULL;

-- Index: filter order berdasarkan status per tenant (untuk order board)
CREATE INDEX IF NOT EXISTS idx_orders_tenant_status ON orders (tenant_id, status) WHERE deleted_at IS NULL;

-- Index: cari order berdasarkan meja (untuk order board per meja)
CREATE INDEX IF NOT EXISTS idx_orders_table_id ON orders (table_id) WHERE deleted_at IS NULL;

-- Index: cari order berdasarkan user (untuk riwayat order customer)
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
