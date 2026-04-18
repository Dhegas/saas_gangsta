BEGIN;

ALTER TABLE subscription_plans
    ADD COLUMN IF NOT EXISTS max_tenants INTEGER;

ALTER TABLE subscription_plans
    ALTER COLUMN max_tenants SET DEFAULT 1;

UPDATE subscription_plans
SET max_tenants = 1
WHERE max_tenants IS NULL;

ALTER TABLE subscription_plans
    DROP CONSTRAINT IF EXISTS subscription_plans_max_tenants_check;

ALTER TABLE subscription_plans
    ADD CONSTRAINT subscription_plans_max_tenants_check
    CHECK (max_tenants IS NULL OR max_tenants > 0);

CREATE TABLE IF NOT EXISTS merchant_tenants (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    is_owner    BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_merchant_tenants_user_tenant UNIQUE (user_id, tenant_id)
);

CREATE INDEX IF NOT EXISTS idx_merchant_tenants_user_id ON merchant_tenants (user_id);
CREATE INDEX IF NOT EXISTS idx_merchant_tenants_tenant_id ON merchant_tenants (tenant_id);

CREATE TRIGGER trg_merchant_tenants_updated_at
    BEFORE UPDATE ON merchant_tenants
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

COMMIT;
