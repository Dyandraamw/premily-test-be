package models

import "time"

type Installment struct {
	Installment_ID string `gorm:"size:100;not null;primary_key"`
	Invoice_ID		string	`gorm:"size:100"`
	Due_Date       time.Time`gorm:"not null"`
	Ins_Amount		float64	`gorm:"type:decimal(16,2);not null"`


}