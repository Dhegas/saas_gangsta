-- Tujuan: Menyimpan detail item pada setiap pesanan sebagai snapshot transaksi.

BEGIN;

CREATE TABLE IF NOT EXISTS order_items (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID          NOT NULL,
    menu_id     UUID          NOT NULL,
    menu_name   VARCHAR(180)  NOT NULL,
    quantity    INTEGER       NOT NULL,
    unit_price  NUMERIC(12,2) NOT NULL,
    subtotal    NUMERIC(12,2) NOT NULL,
    notes       TEXT,
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_order_items_menu
        FOREIGN KEY (menu_id)
        REFERENCES menus(id)
        ON DELETE RESTRICT,
    CONSTRAINT chk_order_items_quantity_positive CHECK (quantity > 0),
    CONSTRAINT chk_order_items_unit_price_non_negative CHECK (unit_price >= 0),
    CONSTRAINT chk_order_items_subtotal_non_negative CHECK (subtotal >= 0)
);

-- Index FK.
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_order_items_menu_id ON order_items (menu_id) WHERE deleted_at IS NULL;

COMMIT;