package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID        string `gorm:"size:100;not null;uniqueIndex;primary_key"`
	Username      string `gorm:"size:255;not null;uniqueIndex"`
	Name          string `gorm:"size:255;not null"`
	Email		  string `gorm:"size:255;not null"`
	Phone         string `gorm:"size:100;not null"`
	Password      string `gorm:"size:255;not null"`
	CompanyName   string `gorm:"size:255;not null"`
	Role          Role   `gorm:"default:staff;not null"`
	RememberToken string `gorm:"size:255;not null"`

	Invoice              []Invoice              `gorm:"foreignKey : UserID"`
	Payment_Status       []Payment_Status       `gorm:"foreignKey: UserID"`
	Statement_Of_Account []Statement_Of_Account `gorm:"foreignKey : UserID"`

	Created_At time.Time
	Updated_At time.Time
}
type Role string

const (
	StaffRole         Role = "staff"
	AdminRole         Role = "admin"
	AccessControlRole Role = "access_control"
)


func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	var err error

	err = db.Debug().Model(User{}).Where("LOWER(email) = ?",strings.ToLower(email)).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) FindByID(db *gorm.DB, userID string) (*User, error) {
	var user User
	var err error

	err = db.Debug().Model(User{}).Where("user_id = ?",userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
