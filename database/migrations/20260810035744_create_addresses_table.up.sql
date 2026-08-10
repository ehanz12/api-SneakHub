CREATE TABLE addresses (
    address_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    nama_penerima VARCHAR(100) NOT NULL,
    nomor_telepon VARCHAR(30) NOT NULL,
    alamat TEXT NOT NULL,
    kota VARCHAR(100) NOT NULL,
    provinsi VARCHAR(100) NOT NULL,
    kode_pos VARCHAR(10) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (address_id),
    KEY idx_addresses_user (user_id),

    CONSTRAINT fk_addresses_user
        FOREIGN KEY (user_id) REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB;
