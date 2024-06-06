package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID        string `gorm:"size:100;not null;uniqueIndex;primary_key"`
	Image         string `gorm:"size:255"`
	Username      string `gorm:"size:255;not null;uniqueIndex"`
	Name          string `gorm:"size:255;not null"`
	Email         string `gorm:"size:255;not null;uniqueIndex"`
	Phone         string `gorm:"size:100;not null"`
	Password      string `gorm:"size:255;not null"`
	CompanyName   string `gorm:"size:255;not null"`
	Role          Role   `gorm:"default:'pending';not null"`
	Verified      Verify `gorm:"default:'pending';not null"`
	RememberToken string `gorm:"size:255;not null"`

	Invoice              []Invoice              `gorm:"foreignKey : UserID"`
	Payment_Status       []Payment_Status       `gorm:"foreignKey: UserID"`
	Statement_Of_Account []Statement_Of_Account `gorm:"foreignKey : UserID"`

	Created_At time.Time
	Updated_At time.Time
}
type Role string
type Verify string

const (
	ActiveVerify  Verify = "active"
	PendingVerify Verify = "pending"
)

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
		Image:       params.Image,
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

func (u *User) UpdateUserPicture(db *gorm.DB, user_id string)  error{
	var user User
	err := db.Debug().First(&user, "user_id = ?", user_id).Error
	if err != nil{
		return fmt.Errorf("Picture not found : %w", err)
	}
	user.Image = u.Image
	user.Updated_At = u.Updated_At

	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to save updated your picture: %w", err)
	}

	return nil
}

func (u *User) GetUnverifiedUser(db *gorm.DB) ([]*User, error) {
	var users []*User
	err := db.Debug().Model(&User{}).Where("verified = ?", PendingVerify).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get unverified users: %v", err)
	}
	return users, nil
}
func (u *User) GetUnroleUser(db *gorm.DB) ([]*User, error) {
	var users []*User
	err := db.Debug().Model(&User{}).Where("role = ?", PendingRole).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get unrole users: %v", err)
	}
	return users, nil
}

func (u *User) VerifyUser(db *gorm.DB, user_id string, verify Verify) error {
	var user User
	if err := db.Debug().Model(User{}).Where("user_id=?", user_id).First(&user).Error; err != nil {
		return err
	}

	user.Verified = verify // Menggunakan nilai verify yang diterima sebagai parameter

	// Update `Updated_At` hanya jika `Verified` berubah menjadi `true`
	if verify == ActiveVerify && user.Verified != ActiveVerify {
		user.Updated_At = time.Now()
	}

	err := db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SetUserRole(db *gorm.DB, user_id string, role Role) error {
	var user User
	if err := db.Debug().Model(User{}).Where("user_id=?", user_id).First(&user).Error; err != nil {
		return err
	}
	user.Role = role
	user.Updated_At = time.Now()

	err := db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
