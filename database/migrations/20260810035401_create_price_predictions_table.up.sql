CREATE TABLE price_predictions (
    prediction_id CHAR(36) NOT NULL,
    seller_id CHAR(36) NOT NULL,
    input_atribut JSON NULL,
    estimated_price_min DECIMAL(15,2) NULL,
    estimated_price_max DECIMAL(15,2) NULL,
    recommended_price DECIMAL(15,2) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (prediction_id),
    KEY idx_price_prediction_seller (seller_id),

    CONSTRAINT fk_price_prediction_seller
        FOREIGN KEY (seller_id) REFERENCES sellers(seller_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
