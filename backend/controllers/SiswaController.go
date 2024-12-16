package controllers

import (
	"log"
	"pendaftaran/database"
	"pendaftaran/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetSiswa(c *fiber.Ctx) error {
	db := database.GetDB()
	var siswaList []models.Siswa

	if err := db.Preload("AsalSekolah").Find(&siswaList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(siswaList)
}

func GetSiswaByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	SiswaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Not Found")
	}

	var Siswa models.Siswa
	if err := db.Preload("AsalSekolah").First(&Siswa, SiswaID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).SendString("not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Siswa)
}

func CreateSiswa(c *fiber.Ctx) error {
	db := database.GetDB()
	var Siswa models.Siswa

	// Parse form data
	if err := c.BodyParser(&Siswa); err != nil {
		log.Println("Error parsing form data:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data: " + err.Error())
	}

	log.Println("Parsed Siswa Data:", Siswa)

	// Handle file upload
	file, err := c.FormFile("foto")
	if err != nil {
		log.Println("Error getting file:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Error getting file: " + err.Error())
	}

	// Save the file
	filepath := "./uploads/" + file.Filename
	if err := c.SaveFile(file, filepath); err != nil {
		log.Println("Error saving file:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving file: " + err.Error())
	}
	Siswa.FotoURL = filepath

	// Validate required fields
	if Siswa.Nama == "" || Siswa.NIK == 0 || Siswa.NISN == 0 || Siswa.Alamat == "" || Siswa.TempatLahir == "" || Siswa.TanggalLahir.IsZero() || Siswa.Gender == "" {
		log.Println("Missing required fields")
		return c.Status(fiber.StatusBadRequest).SendString("All fields are required")
	}

	log.Println("Valid fields passed")

	// Check for existing NIK or NISN
	var existingSiswa models.Siswa
	if err := db.Where("NIK = ? OR NISN = ?", Siswa.NIK, Siswa.NISN).First(&existingSiswa).Error; err == nil {
		log.Println("NIK or NISN already exists")
		return c.Status(fiber.StatusConflict).SendString("NIK or NISN already exists")
	}

	// Create new record
	if err := db.Create(&Siswa).Error; err != nil {
		log.Println("Error creating siswa:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	log.Println("Siswa created successfully")

	return c.JSON(Siswa)
}

func UpdateSiswa(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	SiswaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	var updateSiswa models.Siswa
	// Parse form data
	if err := c.BodyParser(&updateSiswa); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	file, err := c.FormFile("foto")
	if err == nil {
		filepath := "../uploads/" + file.Filename
		if err := c.SaveFile(file, filepath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		updateSiswa.FotoURL = filepath // Only set if a file is uploaded
	}

	if updateSiswa.Nama == "" || updateSiswa.NIK == 0 || updateSiswa.NISN == 0 || updateSiswa.Alamat == "" || updateSiswa.TempatLahir == "" || updateSiswa.TanggalLahir.IsZero() || updateSiswa.Gender == "" {
		return c.Status(fiber.StatusBadRequest).SendString("All fields are required")
	}

	var existingSiswa models.Siswa
	if err := db.Where("NIK = ? OR NISN = ?", updateSiswa.NIK, updateSiswa.NISN).Where("id <> ?", SiswaID).First(&existingSiswa).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("NIK atau NISN sudah ada")
	}

	var Siswa models.Siswa
	if err := db.First(&Siswa, SiswaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not Found Siswa")
	}

	Siswa.Nama = updateSiswa.Nama
	Siswa.NIK = updateSiswa.NIK
	Siswa.NISN = updateSiswa.NISN
	Siswa.Alamat = updateSiswa.Alamat
	Siswa.TempatLahir = updateSiswa.TempatLahir
	Siswa.TanggalLahir = updateSiswa.TanggalLahir
	Siswa.Nilai = updateSiswa.Nilai
	Siswa.JarakTempuh = updateSiswa.JarakTempuh
	Siswa.AsalSekolahID = updateSiswa.AsalSekolahID
	if updateSiswa.FotoURL != "" {
		Siswa.FotoURL = updateSiswa.FotoURL
	}
	if updateSiswa.Gender != "" {
		Siswa.Gender = updateSiswa.Gender
	}
	if err := db.Save(&Siswa).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Siswa)
}

func DeleteSiswa(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	SiswaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := db.Delete(&models.Siswa{}, SiswaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Siswa not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
