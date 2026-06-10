BEGIN;

-- =====================================================
-- Tabel: withdraw_requests
-- Mencatat setiap permintaan withdraw partner
-- amount      = jumlah withdraw yang diminta partner
-- fee_amount  = biaya admin withdraw
-- net_amount  = jumlah yang ditransfer ke partner
-- =====================================================

CREATE TABLE withdraw_requests (
    id              UUID          NOT NULL DEFAULT gen_random_uuid(),
    wallet_id       UUID          NOT NULL,
    user_id         UUID          NOT NULL,

    amount          NUMERIC(15,2) NOT NULL,
    fee_amount      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    net_amount      NUMERIC(15,2) NOT NULL,

    status          VARCHAR(20)   NOT NULL DEFAULT 'PENDING',

    bank_name       VARCHAR(100)  NOT NULL,
    bank_account    VARCHAR(50)   NOT NULL,
    account_holder  VARCHAR(150)  NOT NULL,

    admin_note      TEXT          NULL,

    reviewed_by     UUID          NULL,
    reviewed_at     TIMESTAMPTZ   NULL,
    transferred_at  TIMESTAMPTZ   NULL,

    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT withdraw_requests_pkey PRIMARY KEY (id),

    CONSTRAINT fk_withdraw_requests_wallet
        FOREIGN KEY (wallet_id)
        REFERENCES partner_wallets(id),

    CONSTRAINT fk_withdraw_requests_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT fk_withdraw_requests_reviewed_by
        FOREIGN KEY (reviewed_by)
        REFERENCES users(id),

    CONSTRAINT chk_withdraw_requests_status
        CHECK (
            status IN (
                'PENDING',
                'APPROVED',
                'TRANSFERRED',
                'REJECTED'
            )
        ),

    CONSTRAINT chk_withdraw_requests_amount
        CHECK (amount > 0),

    CONSTRAINT chk_withdraw_requests_fee_amount
        CHECK (fee_amount >= 0),

    CONSTRAINT chk_withdraw_requests_fee_not_greater_than_amount
        CHECK (fee_amount <= amount),

    CONSTRAINT chk_withdraw_requests_net_amount
        CHECK (net_amount >= 0),

    CONSTRAINT chk_withdraw_requests_net_amount_calculation
        CHECK (net_amount = amount - fee_amount)
);

-- =====================================================
-- Indexes
-- =====================================================

CREATE INDEX idx_withdraw_requests_user_id
    ON withdraw_requests (user_id);

CREATE INDEX idx_withdraw_requests_wallet_id
    ON withdraw_requests (wallet_id);

CREATE INDEX idx_withdraw_requests_status
    ON withdraw_requests (status);

CREATE INDEX idx_withdraw_requests_created_at
    ON withdraw_requests (created_at DESC);

COMMIT;