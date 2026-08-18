-- add courier service & Biteship tracking columns to shipments
ALTER TABLE shipments
    ADD COLUMN service VARCHAR(100) NULL COMMENT 'Tipe layanan kurir (REG/ECO) dari Biteship' AFTER kurir,
    ADD COLUMN tracking_id VARCHAR(100) NULL COMMENT 'Waybill / tracking id dari Biteship' AFTER nomor_resi,
    ADD COLUMN tracking_history JSON NULL COMMENT 'Riwayat tracking terakhir dari Biteship' AFTER tracking_id;
