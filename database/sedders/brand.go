package seeders

import (
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/gorm"
)

func SeedBrands(db *gorm.DB) error {
	brands := []string{
		"Nike",
		"Adidas",
		"New Balance",
		"Puma",
		"Converse",
		"Vans",
		"Asics",
		"Reebok",
		"Jordan",
		"Under Armour",
		"Skechers",
		"Salomon",
		"On Running",
		"Hoka",
		"Timberland",
	}

	for _, name := range brands {
		var brand models.Brand

		err := db.
			Where("nama_brand = ?", name).
			First(&brand).Error

		if err == nil {
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		brand = models.Brand{
			NamaBrand: name,
		}

		if err := db.Create(&brand).Error; err != nil {
			return err
		}
	}

	return nil
}
