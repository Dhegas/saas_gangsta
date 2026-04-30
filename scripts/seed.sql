BEGIN;

INSERT INTO tenants (id, name, slug, status)
VALUES ('11111111-1111-1111-1111-111111111111', 'Warung Nusantara', 'warung-nusantara', 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, tenant_id, email, password_hash, role, is_active)
VALUES
	('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', NULL, 'admin@saasgangsta.local', '$2a$10$dummyadminhash', 'ADMIN', TRUE),
	('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', '11111111-1111-1111-1111-111111111111', 'merchant@warung.local', '$2a$10$dummymerchanthash', 'MITRA', TRUE),
	('cccccccc-cccc-cccc-cccc-ccccccccccc1', '11111111-1111-1111-1111-111111111111', 'customer@warung.local', '$2a$10$dummycustomerhash', 'BASIC', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO merchant_profiles (id, tenant_id, store_name, address, phone, opening_hours)
VALUES (
	'22222222-2222-2222-2222-222222222222',
	'11111111-1111-1111-1111-111111111111',
	'Warung Nusantara Cabang 1',
	'Jl. Sudirman No. 1, Jakarta',
	'081234567890',
	'{"monday":"08:00-22:00","tuesday":"08:00-22:00","wednesday":"08:00-22:00","thursday":"08:00-22:00","friday":"08:00-22:00","saturday":"09:00-23:00","sunday":"09:00-23:00"}'::jsonb
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO categories (id, tenant_id, name, description, sort_order, is_active)
VALUES
	('33333333-3333-3333-3333-333333333331', '11111111-1111-1111-1111-111111111111', 'Makanan', 'Menu makanan utama', 1, TRUE),
	('33333333-3333-3333-3333-333333333332', '11111111-1111-1111-1111-111111111111', 'Minuman', 'Pilihan minuman', 2, TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, tenant_id, category_id, name, description, price, is_available)
VALUES
	('44444444-4444-4444-4444-444444444441', '11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333331', 'Nasi Goreng Spesial', 'Nasi goreng dengan telur dan ayam', 28000, TRUE),
	('44444444-4444-4444-4444-444444444442', '11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333331', 'Mie Goreng Jawa', 'Mie goreng manis khas jawa', 26000, TRUE),
	('44444444-4444-4444-4444-444444444443', '11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333332', 'Es Teh Manis', 'Teh manis dingin', 8000, TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO tables (id, tenant_id, table_number, capacity, status)
VALUES
	('55555555-5555-5555-5555-555555555551', '11111111-1111-1111-1111-111111111111', 'A1', 4, 'empty'),
	('55555555-5555-5555-5555-555555555552', '11111111-1111-1111-1111-111111111111', 'A2', 4, 'empty')
ON CONFLICT (id) DO NOTHING;

INSERT INTO subscription_plans (id, name, description, price, billing_cycle, features, is_active)
VALUES
	('66666666-6666-6666-6666-666666666661', 'Basic', 'Paket basic untuk UMKM', 199000, 'monthly', '{"maxMenus":100,"maxUsers":5}'::jsonb, TRUE),
	('66666666-6666-6666-6666-666666666662', 'Pro', 'Paket pro untuk merchant berkembang', 399000, 'monthly', '{"maxMenus":500,"maxUsers":20,"reports":true}'::jsonb, TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO subscriptions (id, tenant_id, plan_id, status, started_at, expires_at)
VALUES (
	'77777777-7777-7777-7777-777777777771',
	'11111111-1111-1111-1111-111111111111',
	'66666666-6666-6666-6666-666666666661',
	'active',
	NOW() - INTERVAL '7 days',
	NOW() + INTERVAL '23 days'
)
ON CONFLICT (id) DO NOTHING;

COMMIT;
