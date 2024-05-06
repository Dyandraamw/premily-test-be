package models

import (
	_ "net/http"
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	Invoice_ID         string    `gorm:"size:100;uniqueIndex;not null;primary_key"`
	UserID             string    `gorm:"size:100"`
	Type               Type      `gorm:"not null"`
	Recipient          string    `gorm:"size:255;not null"`
	Address            string    `gorm:"type:text;not null"`
	Policy_Number      string    `gorm:"size:255;not null"`
	Name_Of_Insured    string    `gorm:"size:255;not null"`
	Address_Of_Insured string    `gorm:"size:255;not null"`
	Type_Of_Insurance  string    `gorm:"size:255;not null"`
	Period_Start       time.Time `gorm:"not null"`
	Period_End         time.Time `gorm:"not null"`
	Terms_Of_Period    string    `gorm:"size:255;not null"`
	Remarks            string    `gorm:"type:text;not null"`

	Description_Details Description_Details   `gorm:"foreignKey:Invoice_ID"`
	Installment         []Installment         `gorm:"foreignKey:Invoice_ID"`
	Sum_Insured_Details []Sum_Insured_Details `gorm:"foreignKey:Invoice_ID"`
	Payment_Status      Payment_Status        `gorm:"foreignKey:Invoice_ID"`

	Created_At time.Time
	Updated_At time.Time
}

type Type string

const (
	CreditType = "credit"
	DebitType  = "debit"
)

func (i *Invoice) GetInvoice(db *gorm.DB) (*[]Invoice, error) {
	var err error
	var invoice []Invoice

	err = db.Debug().Model(&Invoice{}).Find(&invoice).Error
	if err != nil {

		return nil, err
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return &invoice, nil

}
