package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pendaftaran/database"
	"pendaftaran/models"
)

func GetSmp(c *fiber.Ctx) error {
	db := database.GetDB()
	var asalSekolahList []models.AsalSekolah

	if err := db.Find(&asalSekolahList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(asalSekolahList)
}

func GetSmpByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	SmpID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Id")
	}
	var smp models.AsalSekolah
	if err := db.First(&smp, SmpID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(smp)
}

func CreateSmp(c *fiber.Ctx) error {
	db := database.GetDB()
	var AsalSekolah models.AsalSekolah
	if err := c.BodyParser(&AsalSekolah); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if AsalSekolah.Nama == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Nama tidak boleh kosong")
	}

	if err := db.Create(&AsalSekolah).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(AsalSekolah)
}

func UpdateSmp(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	AsalSekolahId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID tidak valid")
	}

	var updateAsalSekolah models.AsalSekolah
	if err := c.BodyParser(&updateAsalSekolah); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if updateAsalSekolah.Nama == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Nama tidak boleh kosong")
	}

	var AsalSekolah models.AsalSekolah
	if err := db.First(&AsalSekolah, AsalSekolahId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Nama Tidak Ada")
	}

	AsalSekolah.Nama = updateAsalSekolah.Nama
	if err := db.Save(&AsalSekolah).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(AsalSekolah)
}

func DeleteSmp(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	AsalSekolahId, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := db.Delete(&models.AsalSekolah{}, AsalSekolahId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Nama Tidak Ada")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
