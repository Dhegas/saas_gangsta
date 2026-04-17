-- ============================================================================
-- 011 — Tabel: subscriptions
-- ============================================================================
-- APA ITU?
--   Tabel subscriptions mencatat STATUS BERLANGGANAN setiap tenant terhadap
--   paket langganan tertentu. Ini menghubungkan tenant dengan subscription_plans.
--
-- FUNGSI:
--   - Mencatat tenant mana yang berlangganan paket mana
--   - Tracking tanggal mulai & berakhir langganan
--   - Status langganan: active, expired, canceled, trial
--   - Admin bisa monitor status semua tenant dari sini
--   - Backend bisa cek apakah tenant masih aktif berlangganan sebelum
--     mengizinkan akses ke fitur tertentu
--
-- CONTOH DATA:
--   | tenant_id | plan_id  | status | started_at          | expires_at          |
--   |-----------|----------|--------|---------------------|---------------------|
--   | uuid-t1   | uuid-pro | active | 2026-04-01 00:00:00 | 2026-05-01 00:00:00 |
--   | uuid-t2   | uuid-bas | trial  | 2026-04-15 00:00:00 | 2026-04-29 00:00:00 |
--   | uuid-t3   | uuid-pro | expired| 2026-03-01 00:00:00 | 2026-04-01 00:00:00 |
--
-- RELASI:
--   subscriptions (N) → (1) tenants             — setiap sub milik satu tenant
--                                                  ON DELETE CASCADE: hapus tenant = hapus sub
--   subscriptions (N) → (1) subscription_plans   — setiap sub merujuk satu paket
--                                                  ON DELETE RESTRICT: jangan hapus paket yang masih dipakai
--
-- CONSTRAINT UNIK:
--   (tenant_id, plan_id, status) — satu tenant tidak boleh punya 2 langganan
--   dengan paket DAN status yang sama secara bersamaan
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS subscriptions (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: langganan ini milik tenant mana
    -- ON DELETE CASCADE: jika tenant dihapus, record langganannya juga dihapus
    tenant_id   UUID         NOT NULL
                REFERENCES tenants(id) ON DELETE CASCADE,

    -- FK ke subscription_plans: tenant berlangganan paket mana
    -- ON DELETE RESTRICT: JANGAN hapus paket yang masih ada subscriber aktifnya
    plan_id     UUID         NOT NULL
                REFERENCES subscription_plans(id) ON DELETE RESTRICT,

    -- Status langganan:
    --   active   = sedang berlangganan & berjalan normal
    --   expired  = masa langganan habis (belum diperpanjang)
    --   canceled = dibatalkan oleh tenant atau admin
    --   trial    = sedang dalam masa trial (percobaan gratis)
    status      VARCHAR(20)  NOT NULL
                CHECK (status IN ('active', 'expired', 'canceled', 'trial')),

    -- Tanggal mulai langganan
    started_at  TIMESTAMPTZ  NOT NULL,

    -- Tanggal berakhir langganan
    -- Setelah melewati tanggal ini, status otomatis jadi expired (via cron/scheduler)
    expires_at  TIMESTAMPTZ  NOT NULL,

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,          -- soft delete

    -- Constraint: satu tenant tidak boleh punya 2 langganan dengan plan & status yang sama
    UNIQUE (tenant_id, plan_id, status)
);

-- Index: cari langganan berdasarkan tenant
CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions (tenant_id) WHERE deleted_at IS NULL;

-- Index: cari langganan berdasarkan paket (untuk admin lihat siapa saja subscriber)
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions (plan_id) WHERE deleted_at IS NULL;

-- Index: filter berdasarkan status (untuk monitoring active/expired)
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions (status) WHERE deleted_at IS NULL;

-- Index: cari langganan yang segera expire (untuk notifikasi/reminder)
CREATE INDEX IF NOT EXISTS idx_subscriptions_expires_at ON subscriptions (expires_at) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
