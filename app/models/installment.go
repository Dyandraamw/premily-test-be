package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Installment struct {
	Installment_ID string `gorm:"size:100;not null;primary_key"`
	Invoice_ID		string	`gorm:"size:100"`
	Due_Date       time.Time`gorm:"not null"`
	Ins_Amount		decimal.Decimal		`gorm:"type:numeric(16,2);not null"`

	Payment_Details		[]Payment_Details	`gorm:"foreign_key:Installment_ID"`

}