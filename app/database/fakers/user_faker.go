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
		UserID: 	uuid.New().String(),
		Username:	faker.Username(),
		Phone:		faker.Phonenumber(),
		Password:	faker.Password(),
		CompanyName:	"PT Nusa Mandiri",
		Role:			"User",
		RememberToken: 	"",
		Created_At:		time.Time{},
		Updated_At:		time.Time{},

	}
}