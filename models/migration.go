package models

import "gorm.io/gorm"

// AutoMigrate is optional. If you already import the provided SQL schema,
// you do not need to run this.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Brand{},
		&Category{},
		&Seller{},
		&SellerTrustScore{},
		&UserActivity{},
		&Product{},
		&ProductImage{},
		&ImageEmbedding{},
		&ConditionScore{},
		&ProductVariant{},
		&RecommendationData{},
		&PriceHistory{},
		&PricePrediction{},
		&MarketPriceData{},
		&Wishlist{},
		&Cart{},
		&CartItem{},
		&Order{},
		&OrderItem{},
		&Payment{},
		&Shipment{},
		&Address{},
		&Review{},
		&Notification{},
	)
}
