-- add full origin address for sellers (required for Biteship booking)
ALTER TABLE sellers
    ADD COLUMN alamat_asal TEXT NULL COMMENT 'Alamat jalan asal pengiriman toko' AFTER kota_asal;
