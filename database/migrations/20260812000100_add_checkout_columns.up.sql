-- add checkout support columns to orders
ALTER TABLE orders
    ADD COLUMN address_id CHAR(36) NULL AFTER seller_id,
    ADD COLUMN subtotal DECIMAL(15,2) NOT NULL DEFAULT 0 AFTER metode_pembayaran,
    ADD COLUMN biaya_pengiriman DECIMAL(15,2) NOT NULL DEFAULT 0 AFTER subtotal,
    ADD KEY idx_orders_address (address_id),
    ADD CONSTRAINT fk_orders_address
        FOREIGN KEY (address_id) REFERENCES addresses(address_id)
        ON UPDATE CASCADE
        ON DELETE SET NULL;

-- payments: allow pending status + store gateway reference & payment url
ALTER TABLE payments
    MODIFY COLUMN status_pembayaran ENUM('pending','paid','failed','expired','refunded') NOT NULL DEFAULT 'pending',
    ADD COLUMN gateway_reference VARCHAR(255) NULL AFTER status_pembayaran,
    ADD COLUMN payment_url VARCHAR(500) NULL AFTER gateway_reference;