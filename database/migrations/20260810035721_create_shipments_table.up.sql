CREATE TABLE shipments (
    shipment_id CHAR(36) NOT NULL,
    order_id CHAR(36) NOT NULL,
    kurir VARCHAR(100) NOT NULL,
    nomor_resi VARCHAR(100) NULL,
    status_pengiriman ENUM(
        'menunggu',
        'dikirim',
        'dalam_perjalanan',
        'sampai'
    ) NOT NULL DEFAULT 'menunggu',
    shipped_at TIMESTAMP NULL,
    delivered_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (shipment_id),
    UNIQUE KEY uq_shipment_order (order_id),
    KEY idx_shipments_resi (nomor_resi),

    CONSTRAINT fk_shipments_order
        FOREIGN KEY (order_id) REFERENCES orders(order_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
