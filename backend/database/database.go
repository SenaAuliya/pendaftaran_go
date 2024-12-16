package database

import (
	"log"
	"pendaftaran/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := "sena:sena123@tcp(127.0.0.1:3306)/pendaftaran-go?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Connection Lost To Database:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Jurusan{}, &models.AsalSekolah{}, &models.Tahun{}, &models.Guru{}, &models.Siswa{}, &models.Pendaftaran{})
	if err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal("Failed to close database connection:", err)
	}
}
