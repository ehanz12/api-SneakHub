package responses

import "time"

type NotificationItemResponse struct {
	NotificationID  string    `json:"notification_id"`
	JenisNotifikasi string    `json:"jenis_notifikasi"`
	IsiNotifikasi   string    `json:"isi_notifikasi"`
	StatusBaca      bool      `json:"status_baca"`
	CreatedAt       time.Time `json:"created_at"`
}

type NotificationListDataResponse struct {
	Items       []NotificationItemResponse `json:"items"`
	UnreadCount int64                      `json:"unread_count"`
}

type NotificationReadResponse struct {
	NotificationID string `json:"notification_id"`
	StatusBaca     bool   `json:"status_baca"`
}
