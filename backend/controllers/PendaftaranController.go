package controllers

import (
	"pendaftaran/database"
	"pendaftaran/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPendaftaran(c *fiber.Ctx) error {
	db := database.GetDB()
	var Pendaftaran []models.Pendaftaran

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	query := c.Query("query", "")
	tipe := c.Query("tipe", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	offset := (page - 1) * limit

	var total int64

	tx := db.Offset(offset).Limit(limit).Preload("Guru").Preload("Siswa").Preload("Tahun").Preload("Jurusan")

	if query != "" {
		tx = tx.Joins("JOIN siswas ON siswas.ID = pendaftarans.siswa_id").Where("siswas.Nama LIKE ?", "%"+query+"%")
	}

	if tipe != "" {
		tx = tx.Where("tipe_pendaftaran = ?", tipe)
	}

	if startDate != "" && endDate != "" {
		tx = tx.Where("tanggal_pendaftaran BETWEEN ? AND ?", startDate, endDate)
	}

	if err := tx.Find(&Pendaftaran).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	countTx := db.Model(&models.Pendaftaran{})
	if query != "" {
		countTx = countTx.Joins("JOIN siswas ON siswas.ID = pendaftarans.siswa_id").Where("siswas.Nama LIKE ?", "%"+query+"%")
	}
	if tipe != "" {
		countTx = countTx.Where("tipe_pendaftaran = ?", tipe)
	}
	if startDate != "" && endDate != "" {
		countTx = countTx.Where("tanggal_pendaftaran BETWEEN ? AND ?", startDate, endDate)
	}
	if err := countTx.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data":       Pendaftaran,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func GetPendaftaranByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	PendaftaranID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}
	var pendaftaran models.Pendaftaran

	if err := db.Preload("Guru").Preload("Siswa").Preload("Tahun").Preload("Jurusan").First(&pendaftaran, PendaftaranID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString((err.Error()))
	}
	return c.JSON(pendaftaran)
}

func CreatePendaftaran(c *fiber.Ctx) error {
	db := database.GetDB()
	var Pendaftaran models.Pendaftaran
	if err := c.BodyParser(&Pendaftaran); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if Pendaftaran.GuruID == 0 || Pendaftaran.SiswaID == 0 || Pendaftaran.TahunID == 0 || Pendaftaran.JurusanID == 0 || Pendaftaran.TanggalPendaftaran.IsZero() || Pendaftaran.TipePendaftaran == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Tidak bisa kosong")
	}

	var guru models.Guru
	if err := db.First(&guru, Pendaftaran.GuruID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Guru not found")
	}

	var siswa models.Siswa
	if err := db.First(&siswa, Pendaftaran.SiswaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Siswa not found")
	}

	var tahun models.Tahun
	if err := db.First(&tahun, Pendaftaran.TahunID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Tahun not found")
	}

	var jurusan models.Jurusan
	if err := db.First(&jurusan, Pendaftaran.JurusanID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Jurusan not found")
	}

	var ExistingPendaftaran models.Pendaftaran
	if err := db.Where("SiswaID = ?", Pendaftaran.SiswaID).First(&ExistingPendaftaran).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("Siswa ID Sudah Ada")
	}

	if err := db.Create(&Pendaftaran).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Pendaftaran)
}

func UpdatePendaftaran(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	PendaftaranID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	var UpdatePendaftaran models.Pendaftaran
	if err := c.BodyParser(&UpdatePendaftaran); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if UpdatePendaftaran.GuruID == 0 || UpdatePendaftaran.SiswaID == 0 || UpdatePendaftaran.TahunID == 0 || UpdatePendaftaran.JurusanID == 0 || UpdatePendaftaran.TanggalPendaftaran.IsZero() || UpdatePendaftaran.TipePendaftaran == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Tidak bisa kosong")
	}

	var guru models.Guru
	if err := db.First(&guru, UpdatePendaftaran.GuruID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Guru not found")
	}

	var siswa models.Siswa
	if err := db.First(&siswa, UpdatePendaftaran.SiswaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Siswa not found")
	}

	var tahun models.Tahun
	if err := db.First(&tahun, UpdatePendaftaran.TahunID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Tahun not found")
	}

	var jurusan models.Jurusan
	if err := db.First(&jurusan, UpdatePendaftaran.JurusanID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Jurusan not found")
	}

	var ExistingPendaftaran models.Pendaftaran
	if err := db.Where("SiswaID = ?", UpdatePendaftaran.SiswaID).First(&ExistingPendaftaran).Error; err != nil {
		return c.Status(fiber.StatusConflict).SendString("Siswa ID Sudah Ada")
	}

	var Pendaftaran models.Pendaftaran
	if err := db.First(&Pendaftaran, PendaftaranID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("NotFound")
	}

	Pendaftaran.GuruID = UpdatePendaftaran.GuruID
	Pendaftaran.SiswaID = UpdatePendaftaran.SiswaID
	Pendaftaran.TahunID = UpdatePendaftaran.TahunID
	Pendaftaran.JurusanID = UpdatePendaftaran.JurusanID
	Pendaftaran.TanggalPendaftaran = UpdatePendaftaran.TanggalPendaftaran
	if UpdatePendaftaran.TipePendaftaran != "" {
		Pendaftaran.TipePendaftaran = UpdatePendaftaran.TipePendaftaran
	}
	if err := db.Save(&Pendaftaran).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Pendaftaran)
}

func DeletePendaftaran(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	PendaftaranID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := db.Delete(&models.Pendaftaran{}, PendaftaranID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Tidak ditemukan")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
