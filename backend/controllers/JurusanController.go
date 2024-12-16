package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"pendaftaran/database"
	"pendaftaran/models"
)

func GetJurusan(c *fiber.Ctx) error {
	db := database.GetDB()
	var jurusanList []models.Jurusan

	if err := db.Find(&jurusanList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(jurusanList)
}

func GetJurusanByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	JurusanID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}
	var jurusan models.Jurusan
	if err := db.First(&jurusan, JurusanID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(jurusan)
}

func CreateJurusan(c *fiber.Ctx) error {
	db := database.GetDB()
	var jurusan models.Jurusan

	// Parse the request body into the jurusan struct
	if err := c.BodyParser(&jurusan); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body: " + err.Error())
	}

	// Validate required fields
	if jurusan.Nama == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Nama tidak boleh kosong")
	}
	if jurusan.KodeJurusan == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Kode Jurusan tidak boleh kosong")
	}

	// Check if Jurusan with the same KodeJurusan exists
	var existingJurusan models.Jurusan
	err := db.Where("kode_jurusan = ?", jurusan.KodeJurusan).First(&existingJurusan).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Handle unexpected errors from the database
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to query database: " + err.Error())
	}

	if existingJurusan.ID != 0 {
		// Record exists, return a conflict error
		return c.Status(fiber.StatusConflict).SendString("Kode Jurusan sudah ada")
	}

	// If no record was found, proceed with creation
	if err := db.Create(&jurusan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create Jurusan: " + err.Error())
	}

	// Return the created Jurusan
	return c.Status(fiber.StatusCreated).JSON(jurusan)
}

func UpdateJurusan(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	jurusanId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID tidak valid")
	}

	var updateJurusan models.Jurusan
	if err := c.BodyParser(&updateJurusan); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if updateJurusan.Nama == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Nama tidak boleh kosong")
	}
	if updateJurusan.KodeJurusan == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Kode Jurusan tidak boleh kosong")
	}

	var jurusan models.Jurusan
	if err := db.First(&jurusan, jurusanId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Jurusan Tidak Ada")
	}

	jurusan.Nama = updateJurusan.Nama
	jurusan.KodeJurusan = updateJurusan.KodeJurusan
	if err := db.Save(&jurusan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(jurusan)
}

func DeleteJurusan(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	jurusanId, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := db.Delete(&models.Jurusan{}, jurusanId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Jurusan Tidak Ada")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
