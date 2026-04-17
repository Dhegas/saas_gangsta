    -- ============================================================================
    -- 000 — Extensions & Utility Functions
    -- ============================================================================
    -- File ini WAJIB dijalankan PERTAMA sebelum semua migration lainnya.
    --
    -- Isi:
    --   1. pgcrypto extension → menyediakan fungsi gen_random_uuid() untuk
    --      generate UUID secara otomatis sebagai primary key.
    --   2. trigger_set_updated_at() → fungsi trigger yang otomatis meng-update
    --      kolom updated_at setiap kali baris di-UPDATE. Tidak perlu manual
    --      set updated_at di kode Go.
    --
    -- Kenapa dipisah?
    --   Karena semua tabel bergantung pada extension dan function ini.
    --   Dengan dipisah, kita pastikan ini selalu ada sebelum tabel dibuat.
    -- ============================================================================

    BEGIN;

    -- Extension pgcrypto: menyediakan gen_random_uuid()
    -- Digunakan oleh semua tabel sebagai DEFAULT value untuk kolom id (UUID).
    CREATE EXTENSION IF NOT EXISTS pgcrypto;

    -- Fungsi trigger: auto-update kolom updated_at
    -- Setiap tabel yang punya kolom updated_at akan menggunakan trigger ini.
    -- Cara kerja: sebelum UPDATE, kolom updated_at otomatis di-set ke NOW().
    CREATE OR REPLACE FUNCTION trigger_set_updated_at()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    COMMIT;
