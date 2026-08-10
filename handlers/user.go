package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/ehanz12/api-SneakHub/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandler(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": "false", "message": "request gagal", "errors": err.Error()})
	}
	if errs := req.Validate(); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "kesalahan validasi",
			"errors":  errs,
		})
	}

	user, err := services.CreateUserService(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "false", "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "false", "message": "gagal membuat token", "errors": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"status": "true", "message": "Registrasi Berhasil", "data": fiber.Map{"user": mappers.ToUserRes(user), "access_token": token}})
}

func LoginHandler(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": "false", "message": "request gagal", "errors": err.Error()})
	}

	if errs := req.ValidateLogin(); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "kesalahan validasi",
			"errors":  errs,
		})
	}

	user, err := services.LoginUserService(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "false", "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "false", "message": "gagal membuat token", "errors": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"status": "true", "message": "Login Berhasil", "data": fiber.Map{"access_token": token, "user": mappers.ToLoginRes(user)}})
}
