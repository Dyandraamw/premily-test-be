package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID        string `gorm:"size:100;not null;uniqueIndex;primary_key"`
	Username      string `gorm:"size:255;not null;uniqueIndex"`
	Name          string `gorm:"size:255;not null"`
	Email         string `gorm:"size:255;not null;uniqueIndex"`
	Phone         string `gorm:"size:100;not null"`
	Password      string `gorm:"size:255;not null"`
	CompanyName   string `gorm:"size:255;not null"`
	Role          Role   `gorm:"default:'pending';not null"`
	Verified      bool   `gorm:"default:false;not null"`
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
	PendingRole       Role = "pending"
)



func (u *User) FindByEmail(db *gorm.DB, email string, password string) (*User, error) {
	var user User
	var err error
	// var password User

	err = db.Debug().Model(User{}).Where("LOWER(email) = ? AND password= ?", strings.ToLower(email), password).First(&user).Error
	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *User) FindEmailRegis(db *gorm.DB, email string) (*User, error) {
	var user User
	var err error
	// var password User

	err = db.Debug().Model(User{}).Where("LOWER(email) = ? ", strings.ToLower(email)).First(&user).Error
	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *User) FindByID(db *gorm.DB, userID string) (*User, error) {
	var user User
	var err error

	err = db.Debug().Model(User{}).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) CreateUser(db *gorm.DB, params *User) (*User, error) {
	user := &User{
		UserID:      params.UserID,
		Username:    params.Username,
		Name:        params.Name,
		Email:       params.Email,
		Phone:       params.Phone,
		Password:    params.Password,
		CompanyName: params.CompanyName,
		Role:        params.Role,
		Verified:    params.Verified,
	}
	err := db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetUnverifiedUser(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Debug().Model(User{}).Where("verified = ?", false).Find(&users).Error
	if err != nil {
		err = fmt.Errorf("Get unverified users fail!: %v", err)
		return nil, err
	}
	return users, nil

}

func (u *User) VerifyAndSetUserRole(db *gorm.DB, user_id string, role Role) error {
	var user User
	if err := db.Debug().Model(User{}).Where("user_id=?", user_id).First(&user).Error; err != nil {
		return err
	}
	user.Verified = true
	user.Role = role
	user.Updated_At = time.Now()

	err := db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
