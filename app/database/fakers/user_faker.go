package fakers

import (
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserFaker(db *gorm.DB) *models.User {
	return &models.User{
		UserID:        uuid.New().String(),
		Username:      faker.Username(),
		Image:         faker.URL() +  "/photos/man-sitting-on-gray-concrete-wall-_M6gy9oHgII/download?force=true",
		Name:          faker.Name(),
		Email:         "budi@gmail.com",
		Phone:         faker.Phonenumber(),
		Password:      faker.Password(),
		CompanyName:   "PT Sinarmas",
		Verified:      "active",
		Role:          "admin",
		RememberToken: "",
		Created_At:    time.Now(),
		Updated_At:    time.Now(),
	}
}
