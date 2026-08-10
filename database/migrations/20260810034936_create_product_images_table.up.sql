CREATE TABLE product_images (
    image_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    url_object_storage TEXT NOT NULL,
    urutan_tampil INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (image_id),
    KEY idx_product_images_product (product_id),

    CONSTRAINT fk_product_images_product
        FOREIGN KEY (product_id) REFERENCES products(product_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
