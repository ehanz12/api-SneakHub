package mappers

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

// displayNotificationType memetakan nilai enum database ke alias Inggris.
func displayNotificationType(jenis string) string {
	switch jenis {
	case "price_alert":
		return "PRICE_ALERT"
	case "restock_alert":
		return "RESTOCK_ALERT"
	case "order_update":
		return "ORDER_UPDATE"
	case "promo":
		return "PROMO"
	case "dll":
		return "DLL"
	}
	return strings.ToUpper(jenis)
}

func ToNotificationListResponse(notifications []models.Notification) []responses.NotificationItemResponse {
	out := make([]responses.NotificationItemResponse, 0, len(notifications))
	for _, n := range notifications {
		out = append(out, responses.NotificationItemResponse{
			NotificationID:  n.NotificationID,
			JenisNotifikasi: displayNotificationType(n.JenisNotifikasi),
			IsiNotifikasi:   n.IsiNotifikasi,
			StatusBaca:      n.StatusBaca,
			CreatedAt:       n.CreatedAt,
		})
	}
	return out
}

func ToNotificationReadResponse(n models.Notification) responses.NotificationReadResponse {
	return responses.NotificationReadResponse{
		NotificationID: n.NotificationID,
		StatusBaca:     n.StatusBaca,
	}
}
