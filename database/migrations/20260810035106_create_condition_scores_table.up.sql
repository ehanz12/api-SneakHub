CREATE TABLE condition_scores (
    condition_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    upper_score DECIMAL(5,2) NULL,
    outsole_score DECIMAL(5,2) NULL,
    midsole_score DECIMAL(5,2) NULL,
    insole_score DECIMAL(5,2) NULL,
    accessories_score DECIMAL(5,2) NULL,
    box_score DECIMAL(5,2) NULL,
    skor_akhir DECIMAL(5,2) NULL,
    dinilai_oleh ENUM('seller', 'admin') NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (condition_id),
    UNIQUE KEY uq_condition_product (product_id),

    CONSTRAINT fk_condition_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
