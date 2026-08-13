package seeders

import (
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) error {
	categories := []string{
		"Running",
		"Sport",
		"Casual",
		"Training",
		"Formal",
		"Boots",
		"Sandal",
		"Limited Edition",
	}

	for _, name := range categories {
		category := models.Category{
			NamaKategori: name,
		}

		if err := db.
			Where("nama_kategori = ?", name).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}
	}

	return nil
}
