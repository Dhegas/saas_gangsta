-- Tujuan: Menambahkan kolom customer_name pada tabel orders untuk mencatat nama pelanggan tanpa registrasi akun (misal POS Kasir / Guest QR).
BEGIN;

ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_name VARCHAR(150);

COMMIT;
