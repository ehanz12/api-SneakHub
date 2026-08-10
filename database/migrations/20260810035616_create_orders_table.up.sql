CREATE TABLE orders (
    order_id CHAR(36) NOT NULL,
    customer_id CHAR(36) NOT NULL,
    seller_id CHAR(36) NOT NULL,
    status_order ENUM(
        'pending',
        'diproses',
        'dikirim',
        'selesai',
        'dibatalkan'
    ) NOT NULL DEFAULT 'pending',
    alamat_pengiriman JSON NOT NULL,
    metode_pembayaran VARCHAR(50) NOT NULL,
    total_pesanan DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (order_id),
    KEY idx_orders_customer (customer_id),
    KEY idx_orders_seller (seller_id),
    KEY idx_orders_status (status_order),

    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_orders_seller
        FOREIGN KEY (seller_id) REFERENCES sellers(seller_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB;
