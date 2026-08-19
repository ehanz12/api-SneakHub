package services

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"gorm.io/gorm"
)

// normalizeKondisi memetakan alias kondisi (mis. NEW) ke nilai enum
// kondisi di database. Mengembalikan "" bila tidak dikenal.
func normalizeKondisi(kondisi string) string {
	switch strings.ToUpper(strings.TrimSpace(kondisi)) {
	case "NEW", "BARU":
		return "new"
	case "USED", "BEKAS", "SECOND":
		return "used"
	case "REFURBISHED", "REKONDISI":
		return "refurbished"
	}
	return ""
}

// resolveBrandIDs mengonversi daftar nama brand menjadi brand_id.
func resolveBrandIDs(names []string) ([]string, error) {
	if len(names) == 0 {
		return nil, nil
	}
	var brands []models.Brand
	if err := database.DB.Select("brand_id").Where("nama_brand IN ?", names).Find(&brands).Error; err != nil {
		return nil, errors.New("gagal memuat data brand")
	}
	ids := make([]string, 0, len(brands))
	for _, b := range brands {
		ids = append(ids, b.BrandID)
	}
	return ids, nil
}

// resolveCategoryIDs mengonversi daftar nama kategori menjadi category_id.
func resolveCategoryIDs(names []string) ([]string, error) {
	if len(names) == 0 {
		return nil, nil
	}
	var categories []models.Category
	if err := database.DB.Select("category_id").Where("nama_kategori IN ?", names).Find(&categories).Error; err != nil {
		return nil, errors.New("gagal memuat data kategori")
	}
	ids := make([]string, 0, len(categories))
	for _, c := range categories {
		ids = append(ids, c.CategoryID)
	}
	return ids, nil
}

// loadUserPreferences mengambil preferensi ukuran dan brand favorit user.
func loadUserPreferences(userID string) (sizes []string, brands []string) {
	if userID == "" {
		return nil, nil
	}
	var user models.User
	if err := database.DB.Select("preferensi_ukuran", "brand_favorit").
		Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, nil
	}
	_ = json.Unmarshal(user.PreferensiUkuran, &sizes)
	_ = json.Unmarshal(user.BrandFavorit, &brands)
	return sizes, brands
}

// saveRecommendation menyimpan daftar rekomendasi user (upsert per user+sumber).
func saveRecommendation(userID, sumber string, productIDs []string) {
	if userID == "" {
		return
	}
	tx := database.DB.Begin()
	if tx.Error != nil {
		return
	}
	if err := tx.Where("user_id = ? AND sumber = ?", userID, sumber).
		Delete(&models.RecommendationData{}).Error; err != nil {
		tx.Rollback()
		return
	}
	if len(productIDs) > 0 {
		data, err := json.Marshal(productIDs)
		if err != nil {
			tx.Rollback()
			return
		}
		if err := tx.Create(&models.RecommendationData{
			UserID:          userID,
			Sumber:          sumber,
			DaftarProductID: data,
		}).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

// SmartFilterService memfilter produk sesuai kriteria user lalu memberi
// match_score berdasarkan prioritas harga, kondisi, dan seller trust.
func SmartFilterService(userID string, r requests.SmartFilterRequest) ([]responses.SmartFilterItemResponse, error) {
	query := database.DB.Model(&models.Product{}).Where("status_publikasi = ?", "aktif")

	if r.BudgetMin > 0 {
		query = query.Where("harga >= ?", r.BudgetMin)
	}
	if r.BudgetMax > 0 {
		query = query.Where("harga <= ?", r.BudgetMax)
	}
	if len(r.Brand) > 0 {
		ids, err := resolveBrandIDs(r.Brand)
		if err != nil {
			return nil, err
		}
		query = query.Where("brand_id IN ?", ids)
	}
	if len(r.Kategori) > 0 {
		ids, err := resolveCategoryIDs(r.Kategori)
		if err != nil {
			return nil, err
		}
		query = query.Where("category_id IN ?", ids)
	}
	if len(r.Kondisi) > 0 {
		kondisis := make([]string, 0, len(r.Kondisi))
		for _, k := range r.Kondisi {
			if kk := normalizeKondisi(k); kk != "" {
				kondisis = append(kondisis, kk)
			}
		}
		if len(kondisis) > 0 {
			query = query.Where("kondisi IN ?", kondisis)
		}
	}
	if len(r.Ukuran) > 0 {
		for _, u := range r.Ukuran {
			query = query.Where("JSON_SEARCH(ukuran_tersedia, 'one', ?) IS NOT NULL", u)
		}
	}

	query = query.Preload("Brand").Preload("Seller").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tampil asc")
	})
	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, errors.New("gagal memuat produk")
	}

	_, favoritBrands := loadUserPreferences(userID)
	favorit := make(map[string]bool)
	for _, b := range favoritBrands {
		favorit[strings.ToLower(b)] = true
	}

	wHarga, wKondisi, wTrust := r.Prioritas.Harga, r.Prioritas.Kondisi, r.Prioritas.SellerTrust
	if wHarga <= 0 && wKondisi <= 0 && wTrust <= 0 {
		wHarga, wKondisi, wTrust = 1, 1, 1
	}
	sumW := wHarga + wKondisi + wTrust

	items := make([]responses.SmartFilterItemResponse, 0, len(products))
	for _, p := range products {
		hargaComp := 100.0
		if r.BudgetMin > 0 && r.BudgetMax > 0 {
			center := (r.BudgetMin + r.BudgetMax) / 2
			half := (r.BudgetMax - r.BudgetMin) / 2
			if half > 0 {
				dist := math.Abs(p.Harga-center) / half
				hargaComp = 100 - dist*50
				if hargaComp < 50 {
					hargaComp = 50
				}
			}
		}

		kondisiComp := 100.0

		trustComp := 50.0
		if p.Seller.SellerTrustScore != nil {
			trustComp = *p.Seller.SellerTrustScore
		}

		match := (wHarga*hargaComp + wKondisi*kondisiComp + wTrust*trustComp) / sumW

		alasan := make([]string, 0, 5)
		if r.BudgetMin > 0 || r.BudgetMax > 0 {
			alasan = append(alasan, "Sesuai budget")
		}
		if len(r.Ukuran) > 0 {
			alasan = append(alasan, "Ukuran tersedia")
		}
		if len(r.Kondisi) > 0 {
			alasan = append(alasan, "Sesuai kondisi")
		}
		if p.Brand.NamaBrand != "" && favorit[strings.ToLower(p.Brand.NamaBrand)] {
			alasan = append(alasan, "Sesuai brand favorit")
		}
		if p.Seller.SellerTrustScore != nil && *p.Seller.SellerTrustScore >= 80 {
			alasan = append(alasan, "Seller terpercaya")
		}

		items = append(items, responses.SmartFilterItemResponse{
			ProductID:  p.ProductID,
			NamaProduk: p.NamaProduk,
			Harga:      p.Harga,
			ImageURL:   FirstImageURL(p.Images),
			MatchScore: int(math.Round(match)),
			Alasan:     alasan,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].MatchScore > items[j].MatchScore })
	if len(items) > 20 {
		items = items[:20]
	}
	return items, nil
}

// ukuranProdukTersedia memeriksa apakah ukuran ada di daftar ukuran produk.
func ukuranProdukTersedia(raw []byte, size string) bool {
	if size == "" {
		return false
	}
	var sizes []string
	if err := json.Unmarshal(raw, &sizes); err != nil {
		return false
	}
	for _, s := range sizes {
		if s == size {
			return true
		}
	}
	return false
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

// PersonalizedRecommendationService merekomendasikan produk berdasarkan
// preferensi (brand & ukuran) serta afinitas dari aktivitas user.
func PersonalizedRecommendationService(userID string, limit int) ([]responses.RecommendationItemResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	sizes, brands := loadUserPreferences(userID)

	var affinityRows []struct {
		BrandID string `gorm:"column:brand_id"`
		Cnt     int64  `gorm:"column:cnt"`
	}
	if err := database.DB.Model(&models.UserActivity{}).
		Select("products.brand_id, COUNT(*) AS cnt").
		Joins("JOIN products ON products.product_id = user_activities.product_id").
		Where("user_activities.user_id = ? AND user_activities.jenis_aktivitas IN ? AND user_activities.product_id IS NOT NULL",
			userID, []string{"view", "add_to_wishlist", "add_to_cart", "purchase"}).
		Group("products.brand_id").
		Scan(&affinityRows).Error; err != nil {
		return nil, errors.New("gagal memuat aktivitas user")
	}
	affinity := make(map[string]int64, len(affinityRows))
	for _, a := range affinityRows {
		affinity[a.BrandID] = a.Cnt
	}

	query := database.DB.Model(&models.Product{}).
		Where("status_publikasi = ?", "aktif").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan_tampil asc")
		}).
		Order("created_at desc").
		Limit(50)

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, errors.New("gagal memuat produk")
	}

	items := make([]responses.RecommendationItemResponse, 0, len(products))
	productIDs := make([]string, 0, len(products))
	for _, p := range products {
		score := 0.1
		brandMatch := len(brands) > 0 && containsString(brands, p.Brand.NamaBrand)
		if brandMatch {
			score += 0.4
		}
		if affinity[p.BrandID] > 0 {
			score += 0.2
		}
		sizeMatch := false
		if len(sizes) > 0 {
			for _, u := range sizes {
				if ukuranProdukTersedia(p.UkuranTersedia, u) {
					sizeMatch = true
					break
				}
			}
		}
		if sizeMatch {
			score += 0.3
		}
		if score > 1 {
			score = 1
		}

		reason := "Rekomendasi berdasarkan aktivitasmu."
		switch {
		case brandMatch && sizeMatch:
			reason = "Sesuai dengan brand dan ukuran favoritmu."
		case brandMatch:
			reason = "Sesuai dengan brand favoritmu."
		case sizeMatch:
			reason = "Sesuai dengan ukuran favoritmu."
		}

		items = append(items, responses.RecommendationItemResponse{
			ProductID:  p.ProductID,
			NamaProduk: p.NamaProduk,
			Harga:      p.Harga,
			ImageURL:   FirstImageURL(p.Images),
			Score:      math.Round(score*100) / 100,
			Reason:     reason,
		})
		if score >= 0.4 {
			productIDs = append(productIDs, p.ProductID)
		}
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	if len(items) > limit {
		items = items[:limit]
	}
	if len(productIDs) > limit {
		productIDs = productIDs[:limit]
	}
	saveRecommendation(userID, "personalized", productIDs)
	return items, nil
}

// TrendingService menghitung skor tren produk berdasarkan views dan
// wishlist dalam periode tertentu.
func TrendingService(userID, period string, limit int) ([]responses.TrendingItemResponse, string, error) {
	if limit <= 0 {
		limit = 10
	}
	period = strings.ToUpper(strings.TrimSpace(period))
	if period == "" {
		period = "WEEKLY"
	}
	switch period {
	case "DAILY", "WEEKLY", "MONTHLY":
	default:
		return nil, "", errors.New("period harus salah satu dari: DAILY, WEEKLY, MONTHLY")
	}

	start := time.Now()
	switch period {
	case "DAILY":
		start = start.AddDate(0, 0, -1)
	case "WEEKLY":
		start = start.AddDate(0, 0, -7)
	case "MONTHLY":
		start = start.AddDate(0, -1, 0)
	}

	viewsMap := make(map[string]int64)
	var views []struct {
		ProductID string `gorm:"column:product_id"`
		Cnt       int64  `gorm:"column:cnt"`
	}
	if err := database.DB.Model(&models.UserActivity{}).
		Select("product_id, COUNT(*) AS cnt").
		Where("jenis_aktivitas = ? AND product_id IS NOT NULL AND created_at >= ?", "view", start).
		Group("product_id").
		Scan(&views).Error; err != nil {
		return nil, "", errors.New("gagal menghitung views")
	}
	for _, v := range views {
		viewsMap[v.ProductID] = v.Cnt
	}

	wishMap := make(map[string]int64)
	var wishlists []struct {
		ProductID string `gorm:"column:product_id"`
		Cnt       int64  `gorm:"column:cnt"`
	}
	if err := database.DB.Model(&models.Wishlist{}).
		Select("product_id, COUNT(*) AS cnt").
		Where("created_at >= ?", start).
		Group("product_id").
		Scan(&wishlists).Error; err != nil {
		return nil, "", errors.New("gagal menghitung wishlist")
	}
	for _, w := range wishlists {
		wishMap[w.ProductID] = w.Cnt
	}

	items := make([]responses.TrendingItemResponse, 0, limit)
	productIDs := make([]string, 0, limit)
	for pid, viewsCount := range viewsMap {
		wishCount := wishMap[pid]
		trend := int(math.Min(100, math.Round(float64(viewsCount)*0.7+float64(wishCount)*1.5)))
		var product models.Product
		if err := database.DB.Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan_tampil asc")
		}).Select("product_id", "nama_produk").
			Where("product_id = ? AND status_publikasi = ?", pid, "aktif").
			First(&product).Error; err != nil {
			continue
		}
		items = append(items, responses.TrendingItemResponse{
			ProductID:     product.ProductID,
			NamaProduk:    product.NamaProduk,
			ImageURL:      FirstImageURL(product.Images),
			TrendScore:    trend,
			Views:         viewsCount,
			WishlistCount: wishCount,
		})
		productIDs = append(productIDs, product.ProductID)
	}

	sort.Slice(items, func(i, j int) bool { return items[i].TrendScore > items[j].TrendScore })
	if len(items) > limit {
		items = items[:limit]
	}
	if len(productIDs) > limit {
		productIDs = productIDs[:limit]
	}
	saveRecommendation(userID, "trending", productIDs)
	return items, period, nil
}

// BestSellerWeeklyService mengambil produk terlaris minggu sebelumnya
// (Senin-Minggu) beserta ranking dan total terjual.
func BestSellerWeeklyService(userID string, limit int) ([]responses.BestSellerItemResponse, time.Time, time.Time, error) {
	if limit <= 0 {
		limit = 10
	}
	now := time.Now()
	weekdayOffset := (int(now.Weekday()) + 6) % 7
	thisMonday := now.AddDate(0, 0, -weekdayOffset)
	thisMonday = time.Date(thisMonday.Year(), thisMonday.Month(), thisMonday.Day(), 0, 0, 0, 0, thisMonday.Location())
	prevMonday := thisMonday.AddDate(0, 0, -7)
	prevSunday := thisMonday.AddDate(0, 0, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var rows []struct {
		ProductID    string `gorm:"column:product_id"`
		NamaProduk   string `gorm:"column:nama_produk"`
		ImageURL     string `gorm:"column:image_url"`
		TotalTerjual int64  `gorm:"column:total_terjual"`
	}
	if err := database.DB.Model(&models.OrderItem{}).
		Select("order_items.product_id, products.nama_produk, (SELECT url_object_storage FROM product_images pi WHERE pi.product_id = order_items.product_id ORDER BY pi.urutan_tampil ASC LIMIT 1) AS image_url, SUM(order_items.jumlah) AS total_terjual").
		Joins("JOIN products ON products.product_id = order_items.product_id").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("orders.created_at >= ? AND orders.created_at <= ? AND orders.status_order <> ?",
			prevMonday, prevSunday, "dibatalkan").
		Group("order_items.product_id, products.nama_produk").
		Order("total_terjual desc").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, time.Time{}, time.Time{}, errors.New("gagal memuat best seller")
	}

	items := make([]responses.BestSellerItemResponse, 0, len(rows))
	productIDs := make([]string, 0, len(rows))
	for i, row := range rows {
		items = append(items, responses.BestSellerItemResponse{
			Rank:         i + 1,
			ProductID:    row.ProductID,
			NamaProduk:   row.NamaProduk,
			ImageURL:     row.ImageURL,
			TotalTerjual: row.TotalTerjual,
		})
		productIDs = append(productIDs, row.ProductID)
	}
	saveRecommendation(userID, "best_seller", productIDs)
	return items, prevMonday, prevSunday, nil
}

// HomePersonalizedService menyusun section homepage: rekomendasi,
// trending, dan best seller mingguan.
func HomePersonalizedService(userID string) ([]responses.HomeSectionResponse, error) {
	recs, err := PersonalizedRecommendationService(userID, 10)
	if err != nil {
		return nil, err
	}
	trends, _, err := TrendingService(userID, "WEEKLY", 10)
	if err != nil {
		return nil, err
	}
	bests, _, _, err := BestSellerWeeklyService(userID, 10)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(recs)+len(trends)+len(bests))
	for _, r := range recs {
		ids = append(ids, r.ProductID)
	}
	for _, t := range trends {
		ids = append(ids, t.ProductID)
	}
	for _, b := range bests {
		ids = append(ids, b.ProductID)
	}

	imageMap := make(map[string]string)
	if len(ids) > 0 {
		var images []models.ProductImage
		if err := database.DB.Where("product_id IN ?", ids).
			Order("urutan_tampil asc").Find(&images).Error; err == nil {
			for _, img := range images {
				if _, ok := imageMap[img.ProductID]; !ok {
					imageMap[img.ProductID] = img.URLObjectStorage
				}
			}
		}
	}

	infoMap := make(map[string]responses.HomeProductResponse)
	for _, r := range recs {
		infoMap[r.ProductID] = responses.HomeProductResponse{ProductID: r.ProductID, NamaProduk: r.NamaProduk, Harga: r.Harga}
	}
	for _, t := range trends {
		if _, ok := infoMap[t.ProductID]; !ok {
			infoMap[t.ProductID] = responses.HomeProductResponse{ProductID: t.ProductID, NamaProduk: t.NamaProduk}
		}
	}
	for _, b := range bests {
		if _, ok := infoMap[b.ProductID]; !ok {
			infoMap[b.ProductID] = responses.HomeProductResponse{ProductID: b.ProductID, NamaProduk: b.NamaProduk}
		}
	}

	buildSection := func(sType, title string, productIDs []string) responses.HomeSectionResponse {
		products := make([]responses.HomeProductResponse, 0, len(productIDs))
		for _, id := range productIDs {
			item := infoMap[id]
			item.ProductID = id
			item.ImageURL = imageMap[id]
			products = append(products, item)
		}
		return responses.HomeSectionResponse{Type: sType, Title: title, Products: products}
	}

	return []responses.HomeSectionResponse{
		buildSection("RECOMMENDED", "Cocok untuk kamu", idsOf(recs)),
		buildSection("TRENDING", "Trending Shoes", idsOfTrends(trends)),
		buildSection("BEST_SELLER", "Best Seller Mingguan", idsOfBests(bests)),
	}, nil
}

func idsOf(recs []responses.RecommendationItemResponse) []string {
	ids := make([]string, 0, len(recs))
	for _, r := range recs {
		ids = append(ids, r.ProductID)
	}
	return ids
}

func idsOfTrends(trends []responses.TrendingItemResponse) []string {
	ids := make([]string, 0, len(trends))
	for _, t := range trends {
		ids = append(ids, t.ProductID)
	}
	return ids
}

func idsOfBests(bests []responses.BestSellerItemResponse) []string {
	ids := make([]string, 0, len(bests))
	for _, b := range bests {
		ids = append(ids, b.ProductID)
	}
	return ids
}
