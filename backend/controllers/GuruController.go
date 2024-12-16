package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pendaftaran/database"
	"pendaftaran/models"
)

func GetGuru(c *fiber.Ctx) error {
	db := database.GetDB()
	var guruList []models.Guru

	if err := db.Find(&guruList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(guruList)
}

func GetGuruByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	GuruID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	var Guru models.Guru
	if err := db.First(&Guru, GuruID).Error; err != nil {
		if err.Error() == " record not found" {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Guru)
}

func CreateGuru(c *fiber.Ctx) error {
	db := database.GetDB()
	var guru models.Guru
	if err := c.BodyParser(&guru); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if guru.Nama == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Nama tidak boleh kosong")
	}

	if err := db.Create(&guru).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(guru)
}

func UpdateGuru(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	guruID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID tidak valid")
	}

	var updateGuru models.Guru
	if err := c.BodyParser(&updateGuru); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if updateGuru.Nama == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Nama tidak boleh kosong")
	}

	var guru models.Guru
	if err := db.First(&guru, guruID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Guru tidak ditemukan")
	}

	guru.Nama = updateGuru.Nama
	if err := db.Save(&guru).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(guru)
}

func DeleteGuru(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	guruID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID tidak valid")
	}

	if err := db.Delete(&models.Guru{}, guruID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Guru tidak ditemukan")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
