-- ============================================================================
-- 003 — Tabel: merchant_profiles
-- ============================================================================
-- APA ITU?
--   Tabel merchant_profiles menyimpan PROFIL TOKO dari setiap merchant/tenant.
--   Ini adalah data yang ditampilkan ke publik (nama toko, alamat, logo, dll).
--
-- FUNGSI:
--   - Menyimpan informasi toko yang bisa diedit oleh merchant
--   - Data ini ditampilkan di halaman toko (customer lihat saat scan QR)
--   - Setiap tenant hanya punya SATU profil merchant (relasi 1:1)
--
-- CONTOH DATA:
--   | store_name       | address            | phone        | opening_hours                    |
--   |------------------|--------------------|--------------|----------------------------------|
--   | Warung Makan Pak | Jl. Merdeka No.10  | 081234567890 | {"mon":"08:00-22:00","tue":"..."} |
--
-- RELASI:
--   merchant_profiles (1) ↔ (1) tenants  — relasi 1:1 (UNIQUE di tenant_id)
--     Artinya: satu tenant hanya punya satu profil merchant, dan sebaliknya.
--     ON DELETE CASCADE: jika tenant dihapus, profil merchant juga ikut dihapus.
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS merchant_profiles (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: UNIQUE memastikan hanya 1 profil per tenant (relasi 1:1)
    -- ON DELETE CASCADE: hapus profil otomatis jika tenant dihapus
    tenant_id       UUID         NOT NULL UNIQUE
                    REFERENCES tenants(id) ON DELETE CASCADE,

    -- Nama toko yang ditampilkan ke customer
    store_name      VARCHAR(180) NOT NULL,

    -- Alamat lengkap toko
    address         TEXT,

    -- Nomor telepon/WhatsApp toko
    phone           VARCHAR(40),

    -- URL logo toko (disimpan di storage seperti Supabase Storage)
    logo_url        TEXT,

    -- Jam buka toko dalam format JSONB
    -- Contoh: {"mon": "08:00-22:00", "tue": "08:00-22:00", "sun": "closed"}
    opening_hours   JSONB        NOT NULL DEFAULT '{}'::jsonb,

    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ           -- soft delete
);

-- Trigger auto-update updated_at
CREATE TRIGGER trg_merchant_profiles_updated_at
    BEFORE UPDATE ON merchant_profiles
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
