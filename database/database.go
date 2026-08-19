package database

import (
	"fmt"
	"log"

	"github.com/ehanz12/api-SneakHub/config"
	seeders "github.com/ehanz12/api-SneakHub/database/seeders"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {

	cfg := config.AppConfig

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FJakarta",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("⚠️ ERROR CONNECT TO DATABASE !", err)
	}

	DB = db

	seeders.SeedBrands(DB)
	seeders.SeedCategories(DB)

	fmt.Println("👌 CONNECT TO DATABASE COMPLETED !")
	return nil
}
