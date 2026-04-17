-- ============================================================================
-- 013 — Tabel: audit_logs
-- ============================================================================
-- APA ITU?
--   Tabel audit_logs menyimpan LOG AKTIVITAS PENTING di platform.
--   Ini adalah tabel IMMUTABLE — data hanya boleh INSERT, TIDAK BOLEH
--   di-UPDATE atau di-DELETE (untuk menjaga integritas audit trail).
--
-- FUNGSI:
--   - Mencatat setiap aksi kritis yang dilakukan user:
--       • ORDER_CANCELED — order dibatalkan
--       • VOID_TRANSACTION — transaksi di-void
--       • MENU_DELETED — menu dihapus
--       • TENANT_SUSPENDED — tenant di-suspend oleh admin
--       • LOGIN_SUCCESS / LOGIN_FAILED — percobaan login
--       • PAYMENT_REFUNDED — pembayaran di-refund
--   - Menyimpan nilai sebelum (old_value) dan sesudah (new_value) perubahan
--   - Menyimpan IP address untuk tracking keamanan
--   - Berguna untuk:
--       • Investigasi masalah / dispute
--       • Compliance audit
--       • Keamanan (deteksi aktivitas mencurigakan)
--
-- CONTOH DATA:
--   | tenant_id | user_id  | action          | entity_type | entity_id | ip_address   |
--   |-----------|----------|-----------------|-------------|-----------|--------------|
--   | uuid-123  | uuid-mrc | ORDER_CANCELED  | order       | uuid-ord  | 192.168.1.10 |
--   | uuid-123  | uuid-mrc | MENU_DELETED    | menu        | uuid-menu | 192.168.1.10 |
--   | NULL      | uuid-adm | TENANT_SUSPENDED| tenant      | uuid-123  | 10.0.0.1     |
--
-- RELASI:
--   audit_logs (N) → (1) tenants  — log ini terkait tenant mana (opsional)
--                                    ON DELETE SET NULL: tenant dihapus → log tetap ada
--   audit_logs (N) → (1) users    — aksi dilakukan oleh user mana (opsional)
--                                    ON DELETE SET NULL: user dihapus → log tetap ada
--
-- CATATAN PENTING:
--   - Tabel ini TIDAK punya updated_at (karena immutable — tidak boleh diubah)
--   - Tabel ini TIDAK punya deleted_at (karena tidak boleh dihapus)
--   - FK menggunakan ON DELETE SET NULL agar log tetap tersimpan
--     meskipun tenant/user yang terkait sudah dihapus
-- ============================================================================

BEGIN;

CREATE TABLE IF NOT EXISTS audit_logs (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),

    -- FK ke tenants: log ini terkait tenant mana
    -- Nullable: beberapa aksi admin bersifat platform-wide
    -- ON DELETE SET NULL: tenant dihapus → log TETAP ADA (data audit wajib dijaga)
    tenant_id   UUID         REFERENCES tenants(id) ON DELETE SET NULL,

    -- FK ke users: siapa yang melakukan aksi ini
    -- Nullable: ada aksi dari system/cron yang tidak punya user
    -- ON DELETE SET NULL: user dihapus → log TETAP ADA
    user_id     UUID         REFERENCES users(id) ON DELETE SET NULL,

    -- Nama aksi yang dilakukan
    -- Contoh: ORDER_CANCELED, MENU_DELETED, VOID_TRANSACTION, LOGIN_FAILED
    action      VARCHAR(80)  NOT NULL,

    -- Nama tipe entity yang terdampak
    -- Contoh: order, menu, payment, tenant, user
    entity_type VARCHAR(80)  NOT NULL,

    -- ID entity yang terdampak (UUID referensi ke tabel terkait)
    entity_id   UUID,

    -- Snapshot data SEBELUM perubahan (opsional, dalam format JSONB)
    -- Contoh: {"status": "pending", "total": 75000}
    old_value   JSONB,

    -- Snapshot data SESUDAH perubahan (opsional, dalam format JSONB)
    -- Contoh: {"status": "canceled", "total": 75000}
    new_value   JSONB,

    -- Data tambahan lain yang relevan (opsional)
    -- Contoh: {"reason": "Customer request", "canceledBy": "merchant"}
    metadata    JSONB        NOT NULL DEFAULT '{}'::jsonb,

    -- IP address yang melakukan aksi (untuk tracking keamanan)
    ip_address  VARCHAR(45),

    -- Hanya created_at — TIDAK ADA updated_at dan deleted_at
    -- Karena audit log bersifat IMMUTABLE (tidak boleh diubah/dihapus)
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Index: filter log per tenant + sort by terbaru
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_created ON audit_logs (tenant_id, created_at DESC);

-- Index: filter log per user
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs (user_id);

-- Index: cari log berdasarkan entity (contoh: semua log untuk order tertentu)
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs (entity_type, entity_id);

-- Index: filter log berdasarkan jenis aksi
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs (action);

COMMIT;
