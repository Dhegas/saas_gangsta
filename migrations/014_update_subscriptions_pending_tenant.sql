BEGIN;

-- Allow subscription onboarding before tenant creation by linking subscription to merchant user first.
ALTER TABLE subscriptions
    ADD COLUMN IF NOT EXISTS subscriber_user_id UUID REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE subscriptions
    ALTER COLUMN tenant_id DROP NOT NULL;

ALTER TABLE subscriptions
    DROP CONSTRAINT IF EXISTS subscriptions_tenant_id_plan_id_status_key;

ALTER TABLE subscriptions
    DROP CONSTRAINT IF EXISTS subscriptions_status_check;

ALTER TABLE subscriptions
    ADD CONSTRAINT subscriptions_status_check
    CHECK (status IN ('active', 'expired', 'canceled', 'trial', 'pending_tenant'));

CREATE INDEX IF NOT EXISTS idx_subscriptions_subscriber_user_id
    ON subscriptions (subscriber_user_id)
    WHERE deleted_at IS NULL;

-- Prevent duplicate pending/active states for the same user-plan before tenant is assigned.
CREATE UNIQUE INDEX IF NOT EXISTS uq_subscriptions_user_plan_status_no_tenant
    ON subscriptions (subscriber_user_id, plan_id, status)
    WHERE tenant_id IS NULL AND deleted_at IS NULL;

-- Preserve uniqueness once a tenant is assigned.
CREATE UNIQUE INDEX IF NOT EXISTS uq_subscriptions_tenant_plan_status
    ON subscriptions (tenant_id, plan_id, status)
    WHERE tenant_id IS NOT NULL AND deleted_at IS NULL;

COMMIT;
