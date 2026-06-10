BEGIN;

-- =====================================================
-- Tabel: wallet_transactions
-- Mencatat seluruh perubahan saldo partner
-------------------------------------------

-- transaction_type:
--   ORDER_PAYMENT   = pendapatan dari order customer
--   WITHDRAW        = partner melakukan withdraw
--   WITHDRAW_REFUND = withdraw ditolak dan saldo dikembalikan
--   ADMIN_ADJUSTMENT= penyesuaian saldo manual oleh admin
----------------------------------------------------------

-- IDEMPOTENCY:
-- Tetap menggunakan unique index (wallet_id, order_id)
-- untuk mencegah satu order menghasilkan kredit saldo lebih dari sekali.
-- =====================================================

CREATE TABLE wallet_transactions (
id                UUID          NOT NULL DEFAULT gen_random_uuid(),
wallet_id         UUID          NOT NULL,
order_id          UUID          NULL,

transaction_type  VARCHAR(30)   NOT NULL,

amount            NUMERIC(15,2) NOT NULL,
fee_amount        NUMERIC(15,2) NOT NULL DEFAULT 0.00,
net_amount        NUMERIC(15,2) NOT NULL,

description       TEXT          NOT NULL DEFAULT '',

created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

CONSTRAINT wallet_transactions_pkey
    PRIMARY KEY (id),

CONSTRAINT fk_wallet_transactions_wallet
    FOREIGN KEY (wallet_id)
    REFERENCES partner_wallets(id),

CONSTRAINT fk_wallet_transactions_order
    FOREIGN KEY (order_id)
    REFERENCES orders(id),

CONSTRAINT chk_wallet_transactions_type
    CHECK (
        transaction_type IN (
            'ORDER_PAYMENT',
            'WITHDRAW',
            'WITHDRAW_REFUND',
            'ADMIN_ADJUSTMENT'
        )
    ),

CONSTRAINT chk_wallet_transactions_amount
    CHECK (amount > 0),

CONSTRAINT chk_wallet_transactions_fee_amount
    CHECK (fee_amount >= 0),

CONSTRAINT chk_wallet_transactions_fee_not_greater_than_amount
    CHECK (fee_amount <= amount),

CONSTRAINT chk_wallet_transactions_net_amount
    CHECK (net_amount >= 0)

);

-- =====================================================
-- Layer 2 Idempotency
-- Satu order tidak boleh menghasilkan transaksi wallet
-- lebih dari satu kali untuk wallet yang sama.
-- =====================================================

CREATE UNIQUE INDEX idx_uq_wallet_credit_per_order
ON wallet_transactions (wallet_id, order_id)
WHERE order_id IS NOT NULL;

-- =====================================================
-- Index performa
-- =====================================================

CREATE INDEX idx_wallet_transactions_wallet_id
ON wallet_transactions (wallet_id);

CREATE INDEX idx_wallet_transactions_order_id
ON wallet_transactions (order_id);

CREATE INDEX idx_wallet_transactions_transaction_type
ON wallet_transactions (transaction_type);

CREATE INDEX idx_wallet_transactions_created_at
ON wallet_transactions (created_at DESC);

COMMIT;
