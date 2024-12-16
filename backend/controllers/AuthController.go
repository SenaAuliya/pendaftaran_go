package controllers

import (
	"strconv"
	"time"

	"pendaftaran/config"
	"pendaftaran/database"
	"pendaftaran/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetUser(c *fiber.Ctx) error {

	db := database.GetDB()

	var admin []models.User

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	query := c.Query("query", "")
	offset := (page - 1) * limit

	var total int64

	if query != "" {
		if err := db.Where("username LIKE ?", "%"+query+"%").Offset(offset).Limit(limit).Find(&admin).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		if err := db.Model(&models.User{}).Where("username LIKE ?", "%"+query+"%").Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if err := db.Offset(offset).Limit(limit).Find(&admin).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		if err := db.Model(&models.User{}).Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.JSON(fiber.Map{
		"data":       admin,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func GetUserByID(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	AdminID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var admin models.User
	if err := db.First(&admin, AdminID).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).SendString("not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(admin)
}

func CreateUser(c *fiber.Ctx) error {
	db := database.GetDB()
	var admin models.User
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if admin.Username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("username tidal boleh kosong")
	}

	if admin.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("password tidak bole kosong ")
	}

	var exisitngAccount models.User
	if err := db.Where("Username = ?", admin.Username).First(&exisitngAccount).Error; err != nil {
		return c.Status(fiber.StatusConflict).SendString("username sudah ada")
	}

	if err := db.Create(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(admin)
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	adminID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("id tidak valid")
	}

	var updateAdmin models.User
	if err := c.BodyParser(&updateAdmin); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if updateAdmin.Username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("username tidal boleh kosong")
	}

	if updateAdmin.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("password tidak bole kosong ")
	}

	var exisitngAccount models.User
	if err := db.Where("Username = ?", updateAdmin.Username).First(&exisitngAccount).Error; err != nil {
		return c.Status(fiber.StatusConflict).SendString("username sudah ada")
	}

	var admin models.User
	if err := db.First(&admin, adminID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("user tidak ada")
	}

	admin.Username = updateAdmin.Username
	admin.Password = updateAdmin.Username
	if err := db.Save(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(admin)
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	adminID, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := db.Delete(&models.User{}, adminID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User Tidak Ada")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func Login(c *fiber.Ctx) error {
	db := database.GetDB()
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := db.Where("username = ? AND password = ?", loginRequest.Username, loginRequest.Password).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "logout successful",
	})
}
