CREATE TABLE notifications (
    notification_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    jenis_notifikasi ENUM(
        'price_alert',
        'restock_alert',
        'order_update',
        'promo',
        'dll'
    ) NOT NULL,
    isi_notifikasi TEXT NOT NULL,
    status_baca BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (notification_id),
    KEY idx_notifications_user (user_id),
    KEY idx_notifications_unread (user_id, status_baca),

    CONSTRAINT fk_notifications_user
        FOREIGN KEY (user_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
