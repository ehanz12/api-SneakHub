CREATE TABLE price_history (
    price_history_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    harga_lama DECIMAL(15,2) NOT NULL,
    harga_baru DECIMAL(15,2) NOT NULL,
    waktu_perubahan TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (price_history_id),
    KEY idx_price_history_product (product_id),

    CONSTRAINT fk_price_history_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
