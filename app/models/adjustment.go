package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Adjustment struct {
	Adjustment_ID     string  `gorm:"size:100;primary_key;not null"`
	Payment_Status_ID string  `gorm:"size:100"`
	Adjustment_Title  string  `gorm:"size:255;not null"`
	Adjustment_Amount Decimal `gorm:"type:numeric(16,2);default:0;not null"`
	Created_At        time.Time
	Updated_At        time.Time
}

func GetTotalWithAdjustments(db *gorm.DB, installmentID string, insAmount Decimal) (Decimal, error) {
	var adjustments []Adjustment
	err := db.Where("installment_id = ?", installmentID).Find(&adjustments).Error
	if err != nil {
		return Decimal{decimal.Zero}, err
	}

	total := insAmount
	for _, adj := range adjustments {
		total = Decimal{total.Add(adj.Adjustment_Amount.Decimal)}
	}
	return total, nil
}
