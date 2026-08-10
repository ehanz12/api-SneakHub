CREATE TABLE reviews (
    review_id CHAR(36) NOT NULL,
    order_id CHAR(36) NOT NULL,
    customer_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    rating DECIMAL(3,2) NOT NULL,
    komentar TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (review_id),
    KEY idx_reviews_order (order_id),
    KEY idx_reviews_customer (customer_id),
    KEY idx_reviews_product (product_id),

    CONSTRAINT fk_reviews_order
        FOREIGN KEY (order_id) REFERENCES orders(order_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_reviews_customer
        FOREIGN KEY (customer_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_reviews_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
