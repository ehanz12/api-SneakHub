-- add origin postal code & city for sellers (shipping cost calculation)
ALTER TABLE sellers
    ADD COLUMN kode_pos_asal VARCHAR(10) NULL COMMENT 'Kode pos asal pengiriman toko' AFTER status_verifikasi,
    ADD COLUMN kota_asal VARCHAR(100) NULL COMMENT 'Kota asal pengiriman toko' AFTER kode_pos_asal;