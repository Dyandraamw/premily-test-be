package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment_Details struct {
	Pay_Detail_ID string 		`gorm:"size:100;primary_key;not null"`
	Installment_ID	string		`gorm:"size:100"`
	Pay_Date      time.Time		`gorm:"not null"`
	Pay_Amount	  decimal.Decimal		`gorm:"type:numeric(16,2);not null"`
	Created_At	  time.Time
	Updated_At	  time.Time
}