-- ============================================================================
-- 012 — Tabel: reviews
-- ============================================================================
-- APA ITU?
--   Tabel reviews menyimpan REVIEW/RATING yang diberikan customer setelah
--   order selesai. Setiap order hanya bisa di-review SATU KALI.
--
-- FUNGSI:
--   - Customer kasih rating (1-5 bintang) dan komentar opsional
--   - Merchant bisa lihat feedback dari customer untuk perbaikan layanan
--   - Data ini bisa dipakai untuk:
--       • Rata-rata rating toko di halaman publik
--       • Laporan kepuasan pelanggan
--       • Menu mana yang paling sering mendapat review baik/buruk
--
-- CONTOH DATA:
--   | tenant_id | order_id | user_id  | rating | comment               |
--   |-----------|----------|----------|--------|-----------------------|
--   | uuid-123  | uuid-ord | uuid-cst | 5      | Makanannya enak!      |
--   | uuid-123  | uuid-or2 | uuid-cs2 | 3      | Agak lama nunggu      |
--
-- RELASI:
--   reviews (N) → (1) tenants  — setiap review milik satu tenant
--                                 ON DELETE CASCADE
--   reviews (1) → (1) orders   — satu review untuk satu order (UNIQUE on order_id)
--                                 ON DELETE CASCADE: hapus order = hapus review
--   reviews (N) → (1) users    — setiap review ditulis oleh satu customer
--                                 ON DELETE CASCADE: hapus user = hapus review-nya
--
-- CATATAN:
--   Endpoint terkait: POST /api/v1/customer/reviews (dari README)
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS reviews (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: review ini untuk toko/tenant mana
    -- ON DELETE CASCADE: hapus tenant = hapus semua review
    tenant_id   UUID         NOT NULL
                REFERENCES tenants(id) ON DELETE CASCADE,

    -- FK ke orders: review ini untuk order mana
    -- UNIQUE: satu order hanya boleh punya satu review (relasi 1:1)
    -- ON DELETE CASCADE: hapus order = hapus review-nya
    order_id    UUID         NOT NULL UNIQUE
                REFERENCES orders(id) ON DELETE CASCADE,

    -- FK ke users: customer yang menulis review
    -- ON DELETE CASCADE: hapus user = hapus review-nya
    user_id     UUID         NOT NULL
                REFERENCES users(id) ON DELETE CASCADE,

    -- Rating bintang (1 sampai 5)
    --   1 = sangat buruk
    --   2 = buruk
    --   3 = biasa
    --   4 = baik
    --   5 = sangat baik
    rating      INTEGER      NOT NULL CHECK (rating BETWEEN 1 AND 5),

    -- Komentar opsional dari customer
    comment     TEXT,

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ           -- soft delete
);

-- Index: filter review per tenant (untuk statistik kepuasan toko)
CREATE INDEX IF NOT EXISTS idx_reviews_tenant_id ON reviews (tenant_id) WHERE deleted_at IS NULL;

-- Index: cari review berdasarkan order
CREATE INDEX IF NOT EXISTS idx_reviews_order_id ON reviews (order_id) WHERE deleted_at IS NULL;

-- Index: cari semua review dari satu customer
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews (user_id) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
