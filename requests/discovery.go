package requests

import "strings"

type PrioritasRequest struct {
	Harga       float64 `json:"harga"`
	Kondisi     float64 `json:"kondisi"`
	SellerTrust float64 `json:"seller_trust"`
}

type SmartFilterRequest struct {
	BudgetMin float64          `json:"budget_min"`
	BudgetMax float64          `json:"budget_max"`
	Brand     []string         `json:"brand"`
	Ukuran    []string         `json:"ukuran"`
	Kondisi   []string         `json:"kondisi"`
	Kategori  []string         `json:"kategori"`
	Prioritas PrioritasRequest `json:"prioritas"`
}

func (r *SmartFilterRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.BudgetMin < 0 {
		errs["budget_min"] = "budget_min tidak boleh negatif"
	}
	if r.BudgetMax < 0 {
		errs["budget_max"] = "budget_max tidak boleh negatif"
	}
	if r.BudgetMin > 0 && r.BudgetMax > 0 && r.BudgetMax < r.BudgetMin {
		errs["budget_max"] = "budget_max harus lebih besar dari budget_min"
	}
	if r.Prioritas.Harga < 0 {
		errs["prioritas.harga"] = "prioritas.harga tidak boleh negatif"
	}
	if r.Prioritas.Kondisi < 0 {
		errs["prioritas.kondisi"] = "prioritas.kondisi tidak boleh negatif"
	}
	if r.Prioritas.SellerTrust < 0 {
		errs["prioritas.seller_trust"] = "prioritas.seller_trust tidak boleh negatif"
	}

	return errs
}

type PricePredictionRequest struct {
	BrandID        string  `json:"brand_id"`
	CategoryID     string  `json:"category_id"`
	Kondisi        string  `json:"kondisi"`
	ConditionScore float64 `json:"condition_score"`
	Ukuran         string  `json:"ukuran"`
	Model          string  `json:"model"`
	TahunRilis     int     `json:"tahun_rilis"`
}

func (r *PricePredictionRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.BrandID) == "" {
		errs["brand_id"] = "brand_id wajib diisi"
	}
	switch strings.ToUpper(strings.TrimSpace(r.Kondisi)) {
	case "NEW", "USED", "REFURBISHED":
	default:
		errs["kondisi"] = "kondisi harus salah satu dari: NEW, USED, REFURBISHED"
	}
	if r.ConditionScore < 0 || r.ConditionScore > 100 {
		errs["condition_score"] = "condition_score harus antara 0 dan 100"
	}
	if r.TahunRilis < 1900 {
		errs["tahun_rilis"] = "tahun_rilis tidak valid"
	}

	return errs
}

type ConditionScoreRequest struct {
	Upper       float64 `json:"upper"`
	Outsole     float64 `json:"outsole"`
	Midsole     float64 `json:"midsole"`
	Insole      float64 `json:"insole"`
	Accessories float64 `json:"accessories"`
	Box         float64 `json:"box"`
	DinilaiOleh string  `json:"dinilai_oleh"`
}

func (r *ConditionScoreRequest) Validate() map[string]string {
	errs := make(map[string]string)

	fields := map[string]float64{
		"upper":       r.Upper,
		"outsole":     r.Outsole,
		"midsole":     r.Midsole,
		"insole":      r.Insole,
		"accessories": r.Accessories,
		"box":         r.Box,
	}
	for name, val := range fields {
		if val < 0 || val > 100 {
			errs[name] = name + " harus antara 0 dan 100"
		}
	}

	switch strings.ToUpper(strings.TrimSpace(r.DinilaiOleh)) {
	case "SELLER", "ADMIN":
	default:
		errs["dinilai_oleh"] = "dinilai_oleh harus salah satu dari: SELLER, ADMIN"
	}

	return errs
}
