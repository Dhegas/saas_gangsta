-- ============================================================================
-- 006 — Tabel: tables
-- ============================================================================
-- APA ITU?
--   Tabel "tables" menyimpan data MEJA FISIK di toko/restoran per tenant.
--   Setiap meja punya nomor unik, kapasitas, dan status real-time.
--
-- FUNGSI:
--   - Merchant bisa kelola daftar meja (CRUD)
--   - Setiap meja punya QR code yang bisa di-scan customer untuk self-order
--   - Status meja otomatis berubah berdasarkan status order:
--       empty    → ada order baru → occupied
--       occupied → order selesai  → empty
--       reserved → dipakai       → occupied
--   - Merchant bisa lihat "order board" per meja di tampilan kasir/POS
--
-- CONTOH DATA:
--   | tenant_id | table_number | capacity | status   | qr_code_url           |
--   |-----------|--------------|----------|----------|-----------------------|
--   | uuid-123  | A1           | 4        | empty    | https://qr.io/tbl-a1  |
--   | uuid-123  | A2           | 2        | occupied | https://qr.io/tbl-a2  |
--   | uuid-123  | VIP-1        | 8        | reserved | https://qr.io/tbl-vip |
--
-- RELASI:
--   tables (N) → (1) tenants  — setiap meja milik satu tenant
--   tables (1) → (N) orders   — satu meja bisa punya banyak order (sepanjang waktu)
--
-- CONSTRAINT UNIK:
--   (tenant_id, table_number) — nomor meja unik per tenant
--   Contoh: Tenant A boleh punya meja "A1", Tenant B juga boleh punya "A1"
--           tapi Tenant A TIDAK boleh punya dua meja bernomor "A1"
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS tables (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: meja ini milik tenant mana
    -- ON DELETE CASCADE: jika tenant dihapus, semua meja-nya juga dihapus
    tenant_id       UUID         NOT NULL
                    REFERENCES tenants(id) ON DELETE CASCADE,

    -- Nomor/nama meja (contoh: "A1", "VIP-1", "Outdoor-3")
    table_number    VARCHAR(50)  NOT NULL,

    -- URL QR code yang mengarah ke halaman self-order untuk meja ini
    -- Customer scan QR → langsung masuk ke menu toko + linked ke meja ini
    qr_code_url     TEXT,

    -- Kapasitas jumlah orang per meja
    -- CHECK: minimal 1 orang
    capacity        INTEGER      NOT NULL CHECK (capacity > 0),

    -- Status real-time meja:
    --   empty    = meja kosong, siap dipakai
    --   occupied = sedang dipakai / ada order aktif
    --   reserved = sudah di-booking / direservasi
    status          VARCHAR(20)  NOT NULL DEFAULT 'empty'
                    CHECK (status IN ('empty', 'occupied', 'reserved')),

    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,          -- soft delete

    -- Constraint: nomor meja unik dalam satu tenant
    UNIQUE (tenant_id, table_number)
);

-- Index: filter meja berdasarkan tenant
CREATE INDEX IF NOT EXISTS idx_tables_tenant_id ON tables (tenant_id) WHERE deleted_at IS NULL;

-- Index: filter meja berdasarkan status per tenant (untuk dashboard meja)
CREATE INDEX IF NOT EXISTS idx_tables_tenant_status ON tables (tenant_id, status) WHERE deleted_at IS NULL;

-- Trigger auto-update updated_at
CREATE TRIGGER trg_tables_updated_at
    BEFORE UPDATE ON tables
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
