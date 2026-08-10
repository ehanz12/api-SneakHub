CREATE TABLE order_items (
    order_item_id CHAR(36) NOT NULL,
    order_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    variant_id CHAR(36) NULL,
    jumlah INT NOT NULL,
    harga_saat_transaksi DECIMAL(15,2) NOT NULL,

    PRIMARY KEY (order_item_id),
    KEY idx_order_items_order (order_id),
    KEY idx_order_items_product (product_id),
    KEY idx_order_items_variant (variant_id),

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id) REFERENCES orders(order_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_order_items_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_order_items_variant
        FOREIGN KEY (variant_id) REFERENCES product_variants(variant_id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
) ENGINE=InnoDB;
