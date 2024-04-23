package models

import "github.com/shopspring/decimal"

type Sum_Insured_Details struct {
	Sum_Insured_ID     string          `gorm:"size:100;not null;primary_key"`
	Invoice_ID         string          `gorm:"size:100"`
	Items_Name         string          `gorm:"size:255;not null"`
	Sum_Insured_Amount decimal.Decimal `gorm:"type:numeric(16,2);not null"`
	Notes              string          `gorm:"size:255;not null"`
}