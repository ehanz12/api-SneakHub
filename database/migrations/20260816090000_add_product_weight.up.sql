-- add weight to products for shipping cost calculation
ALTER TABLE products
    ADD COLUMN berat INT NOT NULL DEFAULT 0 COMMENT 'Berat produk dalam gram' AFTER stok;