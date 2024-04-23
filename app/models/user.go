package models

import (
	"time"

	_ "gorm.io/gorm"
)
type User struct{
	UserID			string `gorm:"size:100;not null;uniqueIndex;primary_key"`
	Username	string `gorm:"size:255;not null;uniqueIndex"`
	Phone		string	`gorm:"size:100;not null"`
	Password	string	`gorm:"size:255;not null"`
	CompanyName	string	`gorm:"size:255;not null"`
	Role 		Role	`gorm:"type:enum('user','admin', 'access_control'); default:user;not null"`
	RememberToken	string `gorm:"size:255;not null"`
	

	Invoice_ID		[]Invoice				`gorm:"foreign_key : UserID"`
	Payment_Status	[]Payment_Status		`gorm:"foreign_key: UserID"`

	Created_At	time.Time
	Updated_At	time.Time
}
type Role string

const(
	UserRole		Role = "user"
	AdminRole		Role = "admin"
	AccessControlRole	Role = "access_control"
)