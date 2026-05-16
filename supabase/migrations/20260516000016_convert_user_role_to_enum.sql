-- Tujuan: Mengonversi kolom role di tabel users menjadi tipe ENUM (PARTNER, CUSTOMER, ADMIN)

BEGIN;

-- 1. Buat tipe enum user_role
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('PARTNER', 'CUSTOMER', 'ADMIN');
    END IF;
END $$;

-- 2. Hapus constraint check yang lama jika ada
DO $$
DECLARE
    r record;
BEGIN
    FOR r IN
        SELECT conname
        FROM pg_constraint
        WHERE conrelid = 'users'::regclass
          AND contype = 'c'
          AND pg_get_constraintdef(oid) ILIKE '%role%'
    LOOP
        EXECUTE format('ALTER TABLE users DROP CONSTRAINT %I', r.conname);
    END LOOP;
END $$;

-- 3. Normalisasi data role ke nilai yang valid untuk ENUM
UPDATE users
SET role = CASE
    WHEN UPPER(TRIM(role)) IN ('BASIC', 'CUSTOMER', 'C') THEN 'CUSTOMER'
    WHEN UPPER(TRIM(role)) IN ('MITRA', 'PARTNER', 'P', 'MERCHANT') THEN 'PARTNER'
    WHEN UPPER(TRIM(role)) IN ('ADMIN', 'A') THEN 'ADMIN'
    ELSE 'CUSTOMER' -- Default jika tidak dikenal
END;

-- 4. Ubah tipe data kolom role menjadi user_role
-- Menggunakan USING role::user_role untuk casting string ke enum
ALTER TABLE users 
    ALTER COLUMN role TYPE user_role 
    USING role::user_role;

COMMIT;
