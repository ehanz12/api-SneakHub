CREATE TABLE recommendation_data (
    recommendation_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    sumber ENUM(
        'personalized',
        'trending',
        'best_seller',
        'dll'
    ) NOT NULL,
    daftar_product_id JSON NOT NULL,
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (recommendation_id),
    KEY idx_recommendation_user (user_id),

    CONSTRAINT fk_recommendation_user
        FOREIGN KEY (user_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
