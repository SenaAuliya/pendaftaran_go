package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pendaftaran/database"
	"pendaftaran/models"
)

func GetTahun(c *fiber.Ctx) error {
	db := database.GetDB()
	var tahunAjaran []models.Tahun

	if err := db.Find(&tahunAjaran).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(tahunAjaran)
}

func GetTahunByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	TahunID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid Id")
	}
	var Tahun models.Tahun
	if err := db.First(&Tahun, TahunID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Tahun)
}

func CreateTahun(c *fiber.Ctx) error {
	db := database.GetDB()
	var Tahun models.Tahun
	if err := c.BodyParser(&Tahun); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if Tahun.Tahun == "" {
		return c.Status(fiber.StatusBadRequest).SendString("tidak boleh kosong")
	}

	if err := db.Create(&Tahun).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Tahun)
}

func UpdateTahun(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	TahunId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID tidak valid")
	}

	var updateTahun models.Tahun
	if err := c.BodyParser(&updateTahun); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if updateTahun.Tahun == "" {
		return c.Status(fiber.StatusBadRequest).SendString("tidak boleh kosong")
	}

	var Tahun models.Tahun
	if err := db.First(&Tahun, TahunId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Tidak Ada")
	}

	Tahun.Tahun = updateTahun.Tahun
	if err := db.Save(&Tahun).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Tahun)
}

func DeleteTahun(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	TahunId, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := db.Delete(&models.Tahun{}, TahunId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Tidak Ada")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
