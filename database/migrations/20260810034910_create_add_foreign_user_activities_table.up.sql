-- Product ID in user activity now can safely reference products.
ALTER TABLE user_activities
    ADD CONSTRAINT fk_activity_product
    FOREIGN KEY (product_id) REFERENCES products(product_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;
