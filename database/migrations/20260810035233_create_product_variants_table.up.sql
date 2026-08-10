CREATE TABLE product_variants (
    variant_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    ukuran VARCHAR(30) NOT NULL,
    stok INT NOT NULL DEFAULT 0,
    harga DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (variant_id),
    UNIQUE KEY uq_product_variant_size (product_id, ukuran),

    CONSTRAINT fk_variants_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
