-- Tambah status moderasi 'pending' ke enum status_publikasi produk.
ALTER TABLE products
    MODIFY COLUMN status_publikasi ENUM('aktif','draft','nonaktif','pending') NOT NULL DEFAULT 'draft';