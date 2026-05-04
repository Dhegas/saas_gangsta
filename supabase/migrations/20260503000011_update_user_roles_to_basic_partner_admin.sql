BEGIN;

LOCK TABLE users IN ACCESS EXCLUSIVE MODE;

ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_role;

UPDATE users 
SET role = 'CUSTOMER' 
WHERE role IS NULL;

UPDATE users
SET role = CASE
    WHEN LOWER(TRIM(role)) IN ('c', 'customer', 'basic') THEN 'CUSTOMER'
    WHEN LOWER(TRIM(role)) IN ('p', 'partner', 'merchant', 'mitra') THEN 'PARTNER'
    WHEN LOWER(TRIM(role)) IN ('a', 'admin') THEN 'ADMIN'
    ELSE UPPER(TRIM(role))
END;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM users
        WHERE role NOT IN ('CUSTOMER', 'PARTNER', 'ADMIN')
    ) THEN
        RAISE EXCEPTION 'Invalid role values exist after normalization';
    END IF;
END $$;

ALTER TABLE users
    ADD CONSTRAINT chk_users_role
    CHECK (role IN ('CUSTOMER', 'PARTNER', 'ADMIN'));

COMMIT;