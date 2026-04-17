-- ============================================================================
-- 010 — Tabel: subscription_plans
-- ============================================================================
-- APA ITU?
--   Tabel subscription_plans menyimpan PAKET BERLANGGANAN yang tersedia
--   di platform SaaS. Ini adalah tabel PLATFORM-LEVEL (tidak terikat tenant).
--   Hanya admin yang bisa mengelola paket ini.
--
-- FUNGSI:
--   - Admin bisa buat berbagai paket langganan (Basic, Pro, Enterprise)
--   - Setiap paket punya harga, siklus billing (bulanan/tahunan), dan fitur
--   - Merchant pilih paket → tercatat di tabel subscriptions
--   - Bisa atur limit per paket (max menu, max meja)
--
-- CONTOH DATA:
--   | name       | price     | billing_cycle | max_menus | max_tables | is_active |
--   |------------|-----------|---------------|-----------|------------|-----------|
--   | Basic      | 99000.00  | monthly       | 20        | 5          | true      |
--   | Pro        | 199000.00 | monthly       | 100       | 20         | true      |
--   | Enterprise | 499000.00 | monthly       | NULL      | NULL       | true      |
--
-- RELASI:
--   subscription_plans (1) → (N) subscriptions  — satu paket bisa dipakai banyak tenant
--
-- CATATAN:
--   Tabel ini TIDAK punya tenant_id karena bersifat platform-level.
--   Hanya admin yang punya akses (role: admin).
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS subscription_plans (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Nama paket (contoh: "Basic", "Pro", "Enterprise")
    name            VARCHAR(120)   NOT NULL,

    -- Deskripsi paket (fitur yang termasuk, benefit, dll)
    description     TEXT,

    -- Harga langganan per siklus billing
    price           NUMERIC(12,2)  NOT NULL CHECK (price >= 0),

    -- Siklus billing:
    --   monthly = bayar per bulan
    --   yearly  = bayar per tahun (biasanya ada diskon)
    billing_cycle   VARCHAR(20)    NOT NULL
                    CHECK (billing_cycle IN ('monthly', 'yearly')),

    -- Daftar fitur dalam format JSONB array
    -- Contoh: ["POS", "Self-Order", "Laporan Harian", "Multi-User"]
    features        JSONB          NOT NULL DEFAULT '[]'::jsonb,

    -- Limit jumlah menu yang boleh dibuat tenant dengan paket ini
    -- NULL = unlimited (tanpa batas)
    max_menus       INTEGER,

    -- Limit jumlah meja yang boleh dibuat tenant dengan paket ini
    -- NULL = unlimited (tanpa batas)
    max_tables      INTEGER,

    -- Toggle aktif/nonaktif paket
    -- Jika false, paket tidak bisa dipilih oleh tenant baru (tapi yang sudah berlangganan tetap jalan)
    is_active       BOOLEAN        NOT NULL DEFAULT TRUE,

    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ             -- soft delete
);

-- Index: filter paket yang aktif (untuk ditampilkan ke merchant saat pilih paket)
CREATE INDEX IF NOT EXISTS idx_subscription_plans_active ON subscription_plans (is_active) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_subscription_plans_updated_at
    BEFORE UPDATE ON subscription_plans
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
