CREATE TABLE image_embeddings (
    embedding_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    image_id CHAR(36) NOT NULL,
    vector JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (embedding_id),
    UNIQUE KEY uq_image_embeddings_image (image_id),
    KEY idx_image_embeddings_product (product_id),

    CONSTRAINT fk_image_embeddings_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_image_embeddings_image
        FOREIGN KEY (image_id) REFERENCES product_images(image_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
