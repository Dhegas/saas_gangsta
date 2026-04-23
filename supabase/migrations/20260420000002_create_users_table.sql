-- Tujuan: Menyimpan akun pengguna lintas peran (BASIC, MITRA, ADMIN).

BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID,
    email          VARCHAR(255) NOT NULL UNIQUE,
    password_hash  TEXT         NOT NULL,
    full_name      VARCHAR(150) NOT NULL,
    role           VARCHAR(20)  NOT NULL CHECK (role IN ('BASIC', 'MITRA', 'ADMIN')),
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    last_login_at  TIMESTAMPTZ,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT fk_users_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE SET NULL
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id) WHERE deleted_at IS NULL;

-- Index email untuk lookup login/user management tanpa row soft-deleted.
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE deleted_at IS NULL;

-- Index kolom yang sering difilter.
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users (is_active) WHERE deleted_at IS NULL;

COMMIT;