package mappers

import (
	"math"
	"strings"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

// displayDinilaiOleh memetakan nilai enum database ke alias Inggris.
func displayDinilaiOleh(v string) string {
	return strings.ToUpper(v)
}

func fvalPtr(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func roundFloat(v float64) float64 {
	return math.Round(v*100) / 100
}

func ToConditionScoreCreateResponse(cs models.ConditionScore) responses.ConditionScoreCreateResponse {
	return responses.ConditionScoreCreateResponse{
		ProductID:   cs.ProductID,
		SkorAkhir:   roundFloat(fvalPtr(cs.SkorAkhir)),
		Upper:       roundFloat(fvalPtr(cs.Upper)),
		Outsole:     roundFloat(fvalPtr(cs.Outsole)),
		Midsole:     roundFloat(fvalPtr(cs.Midsole)),
		Insole:      roundFloat(fvalPtr(cs.Insole)),
		Accessories: roundFloat(fvalPtr(cs.Accessories)),
		Box:         roundFloat(fvalPtr(cs.Box)),
		DinilaiOleh: displayDinilaiOleh(cs.DinilaiOleh),
	}
}

func ToConditionScoreGetResponse(cs models.ConditionScore) responses.ConditionScoreGetResponse {
	return responses.ConditionScoreGetResponse{
		ProductID: cs.ProductID,
		SkorAkhir: roundFloat(fvalPtr(cs.SkorAkhir)),
		Detail: responses.ConditionScoreDetailResponse{
			Upper:       roundFloat(fvalPtr(cs.Upper)),
			Outsole:     roundFloat(fvalPtr(cs.Outsole)),
			Midsole:     roundFloat(fvalPtr(cs.Midsole)),
			Insole:      roundFloat(fvalPtr(cs.Insole)),
			Accessories: roundFloat(fvalPtr(cs.Accessories)),
			Box:         roundFloat(fvalPtr(cs.Box)),
		},
		DinilaiOleh: displayDinilaiOleh(cs.DinilaiOleh),
	}
}
