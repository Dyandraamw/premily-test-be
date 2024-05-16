package models

import (
	"time"

	"gorm.io/gorm"
)

type Statement_Of_Account struct {
	SOA_ID          string `gorm:"size:100;primary_key;not null"`
	UserID         string `gorm:"size:100"`
	Name_Of_Insured string `gorm:"size:255;not null"`
	Period_Start   time.Time	`gorm:"not null"`
	Period_End		time.Time	`gorm:"not null"`

	Statement_Of_Account_Details	[]Statement_Of_Account_Details	`gorm:"foreignKey:SOA_ID"`

	Created_At		time.Time
	Updated_At		time.Time
}

func (s *Statement_Of_Account) CreateNewSOA(db *gorm.DB, soa *Statement_Of_Account)  (*Statement_Of_Account, error) {
	soaModels := &Statement_Of_Account{
		SOA_ID: soa.SOA_ID,
		UserID: soa.UserID,
		Name_Of_Insured: soa.Name_Of_Insured,
		Period_Start: soa.Period_Start,
		Period_End: soa.Period_End,
	}

	err := db.Debug().Create(&soaModels).Error
	if err != nil {
		return nil, err
	}

	return soaModels, nil
}