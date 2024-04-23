package models

import "github.com/shopspring/decimal"

type Description_Details struct {
	Desc_Details_ID string          `gorm:"size:100;not null; primary_key"`
	Invoice_ID      string          `gorm:"size:100;"`
	Premium         decimal.Decimal `gorm:"type:numeric(16,2);not null"`
	Discount        decimal.Decimal `gorm:"type:numeric(16,2);not null"`
	Admin_Cost      decimal.Decimal `gorm:"type:numeric(16,2);not null"`
	Risk_Management decimal.Decimal `gorm:"type:numeric(16,2);"`
	Brokage         decimal.Decimal `gorm:"type:numeric(16,2);"`
	PPH             decimal.Decimal `gorm:"type:numeric(16,2);"`
}