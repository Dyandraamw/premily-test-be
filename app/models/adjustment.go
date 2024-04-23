package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Adjustment struct {
	Adjustment_ID     string  `gorm:"size:100;primary_key;not null"`
	Payment_Status_ID string  `gorm:"size:100"`
	Adjustment_Title  string  `gorm:"size:255;not null"`
	Adjustment_Amount decimal.Decimal		`gorm:"type:numeric(16,2);not null"`
	Created_At        time.Time
	Updated_At			time.Time
}