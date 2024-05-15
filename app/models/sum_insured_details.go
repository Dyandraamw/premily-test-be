package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Sum_Insured_Details struct {
	Sum_Insured_ID     string          `gorm:"size:100;not null;primary_key"`
	Invoice_ID         string          `gorm:"size:100"`
	Items_Name         string          `gorm:"size:255;not null"`
	Sum_Insured_Amount decimal.Decimal `gorm:"type:numeric(16,2);not null"`
	Notes              string          `gorm:"size:255;not null"`
}

func (SumIns *Sum_Insured_Details) CreateSumInsuredDetails(db *gorm.DB, sumInsured *Sum_Insured_Details) (*Sum_Insured_Details, error) {
	sumInsuredModels := &Sum_Insured_Details{
		Sum_Insured_ID:     sumInsured.Sum_Insured_ID,
		Invoice_ID:         sumInsured.Invoice_ID,
		Items_Name:         sumInsured.Items_Name,
		Sum_Insured_Amount: sumInsured.Sum_Insured_Amount,
		Notes:              sumInsured.Notes,
	}

	err := db.Debug().Create(&sumInsuredModels).Error
	if err != nil {
		return nil, err
	}
	return sumInsuredModels, nil
}
