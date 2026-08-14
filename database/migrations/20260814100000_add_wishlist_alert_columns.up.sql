-- Tambah kolom untuk price alert & restock alert pada wishlist.
ALTER TABLE wishlists
    ADD COLUMN price_alert_enabled BOOLEAN NOT NULL DEFAULT FALSE AFTER status_stok,
    ADD COLUMN target_price DECIMAL(15,2) NULL AFTER price_alert_enabled,
    ADD COLUMN restock_alert_enabled BOOLEAN NOT NULL DEFAULT FALSE AFTER target_price;