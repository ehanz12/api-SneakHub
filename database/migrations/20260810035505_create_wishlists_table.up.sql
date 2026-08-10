CREATE TABLE wishlists (
    wishlist_id CHAR(36) NOT NULL,
    customer_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    status_stok ENUM('available', 'out_of_stock') NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (wishlist_id),
    UNIQUE KEY uq_wishlist_customer_product (customer_id, product_id),
    KEY idx_wishlist_product (product_id),

    CONSTRAINT fk_wishlist_customer
        FOREIGN KEY (customer_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_wishlist_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
