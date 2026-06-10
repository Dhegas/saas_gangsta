BEGIN;

-- =====================================================
-- Tabel: partner_wallets
-- Satu partner memiliki tepat satu wallet (UNIQUE user_id)
-- balance TIDAK BOLEH negatif (CHECK constraint)
-- =====================================================

CREATE TABLE partner_wallets (
    id               UUID        NOT NULL DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL,
    balance          NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_earned     NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_withdrawn  NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT partner_wallets_pkey PRIMARY KEY (id),
    CONSTRAINT fk_partner_wallets_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT uq_partner_wallets_user_id UNIQUE (user_id),
    CONSTRAINT chk_partner_wallets_balance CHECK (balance >= 0),
    CONSTRAINT chk_partner_wallets_total_earned CHECK (total_earned >= 0),
    CONSTRAINT chk_partner_wallets_total_withdrawn CHECK (total_withdrawn >= 0)
);

-- Index untuk lookup cepat by user_id
CREATE INDEX idx_partner_wallets_user_id ON partner_wallets (user_id);

COMMIT;
