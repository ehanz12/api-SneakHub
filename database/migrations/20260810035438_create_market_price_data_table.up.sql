CREATE TABLE market_price_data (
    market_price_id CHAR(36) NOT NULL,
    brand_id CHAR(36) NOT NULL,
    category_id CHAR(36) NULL,
    nama_model VARCHAR(200) NOT NULL,
    kondisi ENUM('new', 'used', 'refurbished') NOT NULL,
    ukuran VARCHAR(30) NULL,
    harga DECIMAL(15,2) NOT NULL,
    sumber VARCHAR(255) NULL,
    recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (market_price_id),
    KEY idx_market_brand (brand_id),
    KEY idx_market_category (category_id),

    CONSTRAINT fk_market_price_brand
        FOREIGN KEY (brand_id) REFERENCES brands(brand_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_market_price_category
        FOREIGN KEY (category_id) REFERENCES categories(category_id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
) ENGINE=InnoDB;
