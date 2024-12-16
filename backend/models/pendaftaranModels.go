package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type TipePendaftaran string

const (
	Zonasi   TipePendaftaran = "zonasi"
	Prestasi TipePendaftaran = "prestasi"
)

type Gender string

const (
	Male   Gender = "laki-laki"
	Female Gender = "perempuan"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `gorm:"size:100"`
	Role     string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string
	jwt.RegisteredClaims
}

type Jurusan struct {
	gorm.Model
	Nama        string `gorm:"size:100;not null"`
	KodeJurusan string `gorm:"unique;size:10;not null"`
}

type AsalSekolah struct {
	gorm.Model
	Nama string `gorm:"size:100;not null;unique"`
}

type Tahun struct {
	gorm.Model
	Tahun string `gorm:"size:4;not null"`
}

type Guru struct {
	gorm.Model
	Nama     string `gorm:"size:100;not null"`
	KodeGuru string `gorm:"size:100;unique;not null"`
}

type Siswa struct {
	gorm.Model
	Nama          string      `gorm:"size:100;not null"`
	NIK           int         `gorm:"unique;not null"`
	NISN          int         `gorm:"size:100;unique;not null"`
	AsalSekolahID uint        `gorm:"size:100;not null"`
	AsalSekolah   AsalSekolah `gorm:"foreignKey:AsalSekolahID;references:ID"`
	Nilai         float32     `gorm:"not null"`
	JarakTempuh   float32     `gorm:"not null"`
	Alamat        string      `gorm:"size:255;not null"`
	TempatLahir   string      `gorm:"size:255;not null"`
	TanggalLahir  time.Time   `gorm:"not null"`
	Gender        Gender      `gorm:"type:enum('laki-laki', 'perempuan')"`
	FotoURL       string      `gorm:"size:255"`
}

type Pendaftaran struct {
	gorm.Model
	GuruID             uint            `gorm:"not null"`
	Guru               Guru            `gorm:"foreignKey:GuruID;references:ID"`
	SiswaID            uint            `gorm:"size:100;not null;column:siswa_id"`
	Siswa              Siswa           `gorm:"foreignKey:SiswaID;references:ID"`
	TahunID            uint            `gorm:"not null"`
	Tahun              Tahun           `gorm:"foreignKey:TahunID;references:ID"`
	JurusanID          uint            `gorm:"not null"`
	Jurusan            Jurusan         `gorm:"foreignKey:JurusanID;references:ID"`
	TanggalPendaftaran time.Time       `gorm:"not null"`
	TipePendaftaran    TipePendaftaran `gorm:"type:enum('zonasi', 'prestasi');column:tipe_pendaftaran"`
}
