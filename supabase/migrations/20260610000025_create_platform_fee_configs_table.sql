CREATE TABLE platform_fee_configs (
    id          UUID NOT NULL DEFAULT gen_random_uuid(),
    fee_type    VARCHAR(50) NOT NULL,
    fee_mode    VARCHAR(20) NOT NULL DEFAULT 'FIXED',
    amount      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    is_enabled  BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT platform_fee_configs_pkey PRIMARY KEY (id),
    CONSTRAINT uq_platform_fee_configs_type UNIQUE (fee_type),

    CONSTRAINT chk_fee_type CHECK (
        fee_type IN (
            'TRANSACTION',
            'WITHDRAW',
            'SUBSCRIPTION',
            'DISCOUNT'
        )
    ),

    CONSTRAINT chk_fee_mode CHECK (
        fee_mode IN (
            'FIXED',
            'PERCENTAGE'
        )
    ),

    CONSTRAINT chk_amount CHECK (amount >= 0)
);