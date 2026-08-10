CREATE TABLE cart_items (
    cart_item_id CHAR(36) NOT NULL,
    cart_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    variant_id CHAR(36) NULL,
    jumlah INT NOT NULL DEFAULT 1,
    harga_saat_ditambahkan DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (cart_item_id),
    KEY idx_cart_items_cart (cart_id),
    KEY idx_cart_items_product (product_id),
    KEY idx_cart_items_variant (variant_id),

    CONSTRAINT fk_cart_items_cart
        FOREIGN KEY (cart_id) REFERENCES carts(cart_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_cart_items_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_cart_items_variant
        FOREIGN KEY (variant_id) REFERENCES product_variants(variant_id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
) ENGINE=InnoDB;
