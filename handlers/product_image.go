package handlers

import (
	"io"
	"strconv"

	"github.com/ehanz12/api-SneakHub/config"
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func UploadProductImageHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	fileHeader, err := c.FormFile("gambar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "file gambar wajib diisi", "errors": err.Error()})
	}

	urutanTampil := 0
	if v := c.FormValue("urutan_tampil"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": false, "message": "kesalahan validasi", "errors": fiber.Map{"urutan_tampil": "urutan tampil harus angka >= 0"}})
		}
		urutanTampil = n
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "gagal membaca file gambar", "errors": err.Error()})
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, services.MaxImageSize+1))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "gagal membaca file gambar", "errors": err.Error()})
	}

	baseURL := config.AppConfig.PublicURL
	if baseURL == "" {
		baseURL = c.Protocol() + "://" + c.Hostname()
	}

	image, err := services.UploadProductImageService(userID, productID, data, urutanTampil, baseURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  true,
		"message": "Gambar berhasil diunggah.",
		"data":    mappers.ToProductImageResponse(*image),
	})
}

func DeleteProductImageHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("product_id")
	imageID := c.Params("image_id")

	if err := services.DeleteProductImageService(userID, productID, imageID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Gambar berhasil dihapus.",
		"data":    nil,
	})
}

func ListProductImagesHandler(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	images, err := services.ListProductImagesService(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  true,
		"message": "Daftar gambar produk.",
		"data":    mappers.ToProductImageListResponse(images),
	})
}

func SearchProductByImageHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("gambar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "file gambar wajib diisi", "errors": err.Error()})
	}

	limit := 10
	if v := c.FormValue("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "gagal membaca file gambar", "errors": err.Error()})
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, services.MaxImageSize+1))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "gagal membaca file gambar", "errors": err.Error()})
	}

	results, err := services.SearchProductByImageService(data, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  true,
		"message": "Hasil pencarian berdasarkan gambar.",
		"data":    results,
	})
}
