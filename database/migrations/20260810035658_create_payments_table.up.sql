CREATE TABLE payments (
    payment_id CHAR(36) NOT NULL,
    order_id CHAR(36) NOT NULL,
    metode_pembayaran VARCHAR(50) NOT NULL,
    jumlah DECIMAL(15,2) NOT NULL,
    status_pembayaran ENUM(
        'paid',
        'failed',
        'expired',
        'refunded'
    ) NOT NULL,
    transaction_reference VARCHAR(150) NULL,
    paid_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (payment_id),
    UNIQUE KEY uq_payment_order (order_id),

    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id) REFERENCES orders(order_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
