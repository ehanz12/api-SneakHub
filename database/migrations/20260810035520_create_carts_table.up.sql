CREATE TABLE carts (
    cart_id CHAR(36) NOT NULL,
    customer_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (cart_id),
    UNIQUE KEY uq_cart_customer (customer_id),

    CONSTRAINT fk_cart_customer
        FOREIGN KEY (customer_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
