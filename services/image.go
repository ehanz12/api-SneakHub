package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/deepteams/webp"
	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/utils"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	MaxImageSize   = 5 * 1024 * 1024
	UploadDir      = "uploads"
	defaultLimit   = 10
	maxSearchLimit = 50
)

func UploadProductImageService(userID, productID string, data []byte, urutanTampil int, baseURL string) (*models.ProductImage, error) {
	if len(data) == 0 {
		return nil, errors.New("gambar kosong")
	}
	if len(data) > MaxImageSize {
		return nil, errors.New("ukuran gambar maksimal 5MB")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal terhubung ke database")
	}

	var sellerID models.Seller
	if err := tx.Select("seller_id", "user_id", "status_verifikasi").Where("user_id = ? AND status_verifikasi = ?", userID, "verified").First(&sellerID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user bukan Seller")
	}

	var product models.Product
	if err := tx.Where("product_id = ? AND seller_id = ?", productID, sellerID.SellerID).First(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("product tidak ditemukan")
	}

	img, err := utils.DecodeImage(data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	vector, err := utils.ComputeImageEmbedding(img)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat embedding gambar")
	}
	vectorJSON, err := json.Marshal(vector)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan embedding")
	}

	var webpBuf bytes.Buffer
	if err := webp.Encode(&webpBuf, img, &webp.EncoderOptions{Quality: 80, Method: 4}); err != nil {
		tx.Rollback()
		return nil, errors.New("gagal mengonversi gambar ke WebP")
	}

	imageID := uuid.NewString()
	fileName := imageID + ".webp"
	filePath := filepath.Join(UploadDir, fileName)

	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyiapkan folder uploads")
	}
	if err := os.WriteFile(filePath, webpBuf.Bytes(), 0644); err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan file gambar")
	}

	image := models.ProductImage{
		ImageID:          imageID,
		ProductID:        productID,
		URLObjectStorage: baseURL + "/uploads/" + fileName,
		UrutanTampil:     urutanTampil,
	}
	if err := tx.Create(&image).Error; err != nil {
		os.Remove(filePath)
		tx.Rollback()
		return nil, errors.New("gagal menyimpan data gambar")
	}

	embedding := models.ImageEmbedding{
		ProductID: productID,
		ImageID:   imageID,
		Vector:    datatypes.JSON(vectorJSON),
	}
	if err := tx.Create(&embedding).Error; err != nil {
		os.Remove(filePath)
		tx.Rollback()
		return nil, errors.New("gagal menyimpan embedding gambar")
	}

	if err := tx.Commit().Error; err != nil {
		os.Remove(filePath)
		tx.Rollback()
		return nil, errors.New("data gambar gagal tersimpan")
	}

	return &image, nil
}

func DeleteProductImageService(userID, productID, imageID string) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("gagal terhubung ke database")
	}

	var sellerID models.Seller
	if err := tx.Select("seller_id", "user_id", "status_verifikasi").Where("user_id = ? AND status_verifikasi = ?", userID, "verified").First(&sellerID).Error; err != nil {
		tx.Rollback()
		return errors.New("user bukan Seller")
	}

	var product models.Product
	if err := tx.Where("product_id = ? AND seller_id = ?", productID, sellerID.SellerID).First(&product).Error; err != nil {
		tx.Rollback()
		return errors.New("product tidak ditemukan")
	}

	var image models.ProductImage
	if err := tx.Where("image_id = ? AND product_id = ?", imageID, productID).First(&image).Error; err != nil {
		tx.Rollback()
		return errors.New("gambar tidak ditemukan")
	}

	if err := tx.Where("image_id = ?", imageID).Delete(&models.ImageEmbedding{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus embedding gambar")
	}

	if err := tx.Delete(&image).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data gambar")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data")
	}

	if fileName := filepath.Base(image.URLObjectStorage); fileName != "." && fileName != "/" {
		os.Remove(filepath.Join(UploadDir, fileName))
	}

	return nil
}

func ListProductImagesService(productID string) ([]models.ProductImage, error) {
	var count int64
	if err := database.DB.Model(&models.Product{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return nil, errors.New("gagal memuat produk")
	}
	if count == 0 {
		return nil, errors.New("product tidak ditemukan")
	}

	var images []models.ProductImage
	if err := database.DB.Preload("Embedding").Where("product_id = ?", productID).Order("urutan_tampil ASC").Find(&images).Error; err != nil {
		return nil, errors.New("gagal memuat gambar produk")
	}
	return images, nil
}

func SearchProductByImageService(data []byte, limit int) ([]responses.SearchProductByImageResponse, error) {
	if len(data) == 0 {
		return nil, errors.New("gambar kosong")
	}
	if len(data) > MaxImageSize {
		return nil, errors.New("ukuran gambar maksimal 5MB")
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxSearchLimit {
		limit = maxSearchLimit
	}

	img, err := utils.DecodeImage(data)
	if err != nil {
		return nil, err
	}
	queryVector, err := utils.ComputeImageEmbedding(img)
	if err != nil {
		return nil, errors.New("gagal membuat embedding gambar")
	}

	var embeddings []models.ImageEmbedding
	if err := database.DB.Preload("Product").Preload("Image").Find(&embeddings).Error; err != nil {
		return nil, errors.New("gagal memuat data embedding")
	}
	if len(embeddings) == 0 {
		return nil, errors.New("belum ada gambar terdaftar untuk dicari")
	}

	type scored struct {
		score float64
		resp  responses.SearchProductByImageResponse
	}
	results := make([]scored, 0, len(embeddings))
	for _, e := range embeddings {
		var vec []float64
		if err := json.Unmarshal(e.Vector, &vec); err != nil || len(vec) == 0 {
			continue
		}
		score := utils.CosineSimilarity(queryVector, vec)
		results = append(results, scored{
			score: score,
			resp: responses.SearchProductByImageResponse{
				ProductID:     e.Product.ProductID,
				NamaProduk:    e.Product.NamaProduk,
				Harga:         e.Product.Harga,
				ImageURL:      e.Image.URLObjectStorage,
				SkorKemiripan: math.Round(score*10000) / 10000,
			},
		})
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	if len(results) > limit {
		results = results[:limit]
	}

	out := make([]responses.SearchProductByImageResponse, len(results))
	for i := range results {
		out[i] = results[i].resp
	}
	return out, nil
}
