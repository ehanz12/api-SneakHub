package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToAddressResponse(a models.Address) responses.AddressResponse {
	return responses.AddressResponse{
		AddressID:    a.AddressID,
		NamaPenerima: a.NamaPenerima,
		NomorTelepon: a.NomorTelepon,
		Alamat:       a.Alamat,
		Kota:         a.Kota,
		Provinsi:     a.Provinsi,
		KodePos:      a.KodePos,
		IsDefault:    a.IsDefault,
		CreatedAt:    a.CreatedAt,
	}
}

func ToAddressListResponse(addresses []models.Address) []responses.AddressResponse {
	out := make([]responses.AddressResponse, 0, len(addresses))
	for _, a := range addresses {
		out = append(out, ToAddressResponse(a))
	}
	return out
}