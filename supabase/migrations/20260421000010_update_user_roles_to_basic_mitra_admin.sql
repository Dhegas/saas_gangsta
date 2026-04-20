-- Tujuan: Standarisasi nilai role users menjadi BASIC, MITRA, ADMIN.

BEGIN;

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

UPDATE users
SET role = CASE
    WHEN LOWER(TRIM(role)) IN ('c', 'customer', 'basic') THEN 'BASIC'
    WHEN LOWER(TRIM(role)) IN ('m', 'merchant', 'mitra') THEN 'MITRA'
    WHEN LOWER(TRIM(role)) IN ('a', 'admin') THEN 'ADMIN'
    ELSE UPPER(TRIM(role))
END;

ALTER TABLE users
    ADD CONSTRAINT chk_users_role
    CHECK (role IN ('BASIC', 'MITRA', 'ADMIN'));

COMMIT;
