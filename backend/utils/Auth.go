package utils

import (
	"errors"
	"pendaftaran/models"

	"gorm.io/gorm"
)

func FinByCredetials(db *gorm.DB, username, password string) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			return nil, errors.New("user not found or incorrect password")
		}
		return nil, err
	}
	return &user, nil
}
