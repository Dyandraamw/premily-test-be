package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Statement_Of_Account struct {
	SOA_ID          string    `gorm:"size:100;primary_key;not null"`
	UserID          string    `gorm:"size:100"`
	Name_Of_Insured string    `gorm:"size:255;not null"`
	Period_Start    time.Time `gorm:"not null;default:current_timestamp"`
	Period_End      time.Time `gorm:"not null;default:current_timestamp"`

	Statement_Of_Account_Details []Statement_Of_Account_Details `gorm:"foreignKey:SOA_ID;constraint:OnDelete:CASCADE"`

	Created_At time.Time
	Updated_At time.Time
}

func (s *Statement_Of_Account) CreateNewSOA(db *gorm.DB, soa *Statement_Of_Account) (*Statement_Of_Account, error) {
	soaModels := &Statement_Of_Account{
		SOA_ID:          soa.SOA_ID,
		UserID:          soa.UserID,
		Name_Of_Insured: soa.Name_Of_Insured,
		Period_Start:    soa.Period_Start,
		Period_End:      soa.Period_End,
		Created_At:      time.Now(),
		Updated_At:      time.Now(),
	}

	err := db.Debug().Create(&soaModels).Error
	if err != nil {
		return nil, err
	}

	return soaModels, nil
}

func (s *Statement_Of_Account) DeleteSOA(db *gorm.DB, soa_id string) error {
	soa := &Statement_Of_Account{}
	if err := db.Debug().First(&soa, "soa_id = ? ", soa_id).Error; err != nil {
		return err
	}
	if err := db.Delete(&soa).Error; err != nil {
		fmt.Printf("Fail!")
		return err
	}
	return nil
}
