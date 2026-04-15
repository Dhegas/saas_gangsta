BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS tenants (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(150) NOT NULL,
	slug VARCHAR(80) NOT NULL UNIQUE,
	status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID REFERENCES tenants(id),
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	role VARCHAR(20) NOT NULL CHECK (role IN ('customer', 'merchant', 'admin')),
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

CREATE TABLE IF NOT EXISTS merchant_profiles (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(id),
	store_name VARCHAR(180) NOT NULL,
	address TEXT,
	phone VARCHAR(40),
	logo_url TEXT,
	opening_hours JSONB NOT NULL DEFAULT '{}'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID NOT NULL REFERENCES tenants(id),
	name VARCHAR(120) NOT NULL,
	description TEXT,
	sort_order INTEGER NOT NULL DEFAULT 0,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (tenant_id, name)
);

CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories (tenant_id);

CREATE TABLE IF NOT EXISTS menus (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID NOT NULL REFERENCES tenants(id),
	category_id UUID REFERENCES categories(id),
	name VARCHAR(180) NOT NULL,
	description TEXT,
	price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
	image_url TEXT,
	is_available BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_menus_tenant_id ON menus (tenant_id);
CREATE INDEX IF NOT EXISTS idx_menus_category_id ON menus (category_id);
CREATE INDEX IF NOT EXISTS idx_menus_tenant_available ON menus (tenant_id, is_available);

CREATE TABLE IF NOT EXISTS tables (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID NOT NULL REFERENCES tenants(id),
	table_number VARCHAR(50) NOT NULL,
	capacity INTEGER NOT NULL CHECK (capacity > 0),
	status VARCHAR(20) NOT NULL DEFAULT 'empty' CHECK (status IN ('empty', 'occupied', 'reserved')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (tenant_id, table_number)
);

CREATE INDEX IF NOT EXISTS idx_tables_tenant_id ON tables (tenant_id);

CREATE TABLE IF NOT EXISTS orders (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID NOT NULL REFERENCES tenants(id),
	table_id UUID REFERENCES tables(id),
	user_id UUID REFERENCES users(id),
	idempotency_key VARCHAR(120) UNIQUE,
	status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'cooking', 'ready', 'done', 'canceled')),
	subtotal NUMERIC(12,2) NOT NULL DEFAULT 0,
	tax NUMERIC(12,2) NOT NULL DEFAULT 0,
	total NUMERIC(12,2) NOT NULL DEFAULT 0,
	notes TEXT,
	order_source VARCHAR(20) NOT NULL CHECK (order_source IN ('self_order', 'pos')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_tenant_id_created_at ON orders (tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_tenant_status ON orders (tenant_id, status);

CREATE TABLE IF NOT EXISTS order_items (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
	menu_id UUID NOT NULL REFERENCES menus(id),
	quantity INTEGER NOT NULL CHECK (quantity > 0),
	unit_price NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
	subtotal NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),
	notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);

CREATE TABLE IF NOT EXISTS payments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	order_id UUID NOT NULL REFERENCES orders(id),
	tenant_id UUID NOT NULL REFERENCES tenants(id),
	idempotency_key VARCHAR(120) UNIQUE,
	method VARCHAR(20) NOT NULL CHECK (method IN ('cash', 'qris', 'transfer')),
	status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed', 'refunded')),
	amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
	paid_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_tenant_id_created_at ON payments (tenant_id, created_at DESC);

CREATE TABLE IF NOT EXISTS subscription_plans (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(120) NOT NULL,
	description TEXT,
	price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
	billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('monthly', 'yearly')),
	features JSONB NOT NULL DEFAULT '{}'::jsonb,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscriptions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID NOT NULL REFERENCES tenants(id),
	plan_id UUID NOT NULL REFERENCES subscription_plans(id),
	status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'expired', 'canceled', 'trial')),
	started_at TIMESTAMPTZ NOT NULL,
	expires_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (tenant_id, plan_id, status)
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions (tenant_id);

CREATE TABLE IF NOT EXISTS audit_logs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id UUID REFERENCES tenants(id),
	user_id UUID REFERENCES users(id),
	action VARCHAR(80) NOT NULL,
	entity_type VARCHAR(80) NOT NULL,
	entity_id UUID,
	metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id_created_at ON audit_logs (tenant_id, created_at DESC);

COMMIT;
