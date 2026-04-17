-- ============================================================================
-- 002 — Tabel: users
-- ============================================================================
-- APA ITU?
--   Tabel users menyimpan data SEMUA pengguna platform, baik itu customer
--   (pelanggan toko), merchant (pemilik toko), maupun admin (pengelola SaaS).
--
-- FUNGSI:
--   - Autentikasi: menyimpan email & password hash untuk login
--   - Otorisasi: kolom "role" menentukan akses (customer/merchant/admin)
--   - Multi-tenant link: merchant & customer terhubung ke tenant,
--     sedangkan admin TIDAK (tenant_id nullable)
--
-- CONTOH DATA:
--   | email              | role     | tenant_id | is_active |
--   |--------------------|----------|-----------|-----------|
--   | john@toko.com      | merchant | uuid-123  | true      |
--   | budi@gmail.com     | customer | uuid-123  | true      |
--   | superadmin@saas.id | admin    | NULL      | true      |
--
-- RELASI:
--   users (N) → (1) tenants        — setiap user terikat ke satu tenant
--                                     (kecuali admin → tenant_id = NULL)
--   users (1) → (N) orders         — satu user bisa punya banyak order
--   users (1) → (N) reviews        — satu user bisa kasih banyak review
--   users (1) → (N) audit_logs     — aktivitas user dicatat di audit log
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS users (
    -- Primary key UUID
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: menghubungkan user dengan toko/tenant tertentu
    -- Nullable karena admin platform tidak terikat tenant manapun
    -- ON DELETE SET NULL: jika tenant dihapus, user tetap ada tapi tenant_id jadi NULL
    tenant_id       UUID         REFERENCES tenants(id) ON DELETE SET NULL,

    -- Email unik untuk login (tidak boleh duplikat di seluruh platform)
    email           VARCHAR(255) NOT NULL UNIQUE,

    -- Password yang sudah di-hash (bcrypt/argon2), JANGAN simpan plain text!
    password_hash   TEXT         NOT NULL,

    -- Nama lengkap user (untuk display di UI)
    full_name       VARCHAR(150),

    -- Role menentukan level akses:
    --   customer = pelanggan, hanya bisa lihat menu & order
    --   merchant = pemilik toko, bisa kelola menu/order/transaksi
    --   admin    = pengelola platform, bisa kelola semua tenant
    role            VARCHAR(20)  NOT NULL
                    CHECK (role IN ('customer', 'merchant', 'admin')),

    -- Flag aktif/non-aktif (bisa di-toggle oleh admin)
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,

    -- Timestamp terakhir kali user login (untuk tracking/analytics)
    last_login_at   TIMESTAMPTZ,

    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ           -- soft delete
);

-- Index: filter user berdasarkan tenant (sering dipakai di query merchant/customer)
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id) WHERE deleted_at IS NULL;

-- Index: filter user berdasarkan role (untuk admin dashboard)
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role) WHERE deleted_at IS NULL;

-- Index: lookup user berdasarkan email (dipakai saat login)
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
