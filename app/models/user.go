package models

import (
	"time"

	_ "gorm.io/gorm"
)
type User struct{
	UserID			string `gorm:"size:100;not null;uniqueIndex;primary_key"`
	Username	string `gorm:"size:255;not null;uniqueIndex"`
	Phone		string	`gorm:"not null"`
	Password	string	`gorm:"size:255;not null"`
	CompanyName	string	`gorm:"size:255;not null"`
	Role 		Role	`gorm:"type:enum('user','admin', 'access_control'); default:user;not null"`
	RememberToken	string `gorm:"size:255;not null"`
	Created_at	time.Time
	Update_At	time.Time
}
type Role string

const(
	UserRole		Role = "user"
	AdminRole		Role = "admin"
	AccessControlRole	Role = "access_control"
)