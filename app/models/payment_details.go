package models

import "time"

type Payment_Details struct {
	Pay_Detail_ID string 		`gorm:"size:100;primary_key;not null"`
	Installment_ID	string		`gorm:"size:100"`
	Pay_Date      time.Time		`gorm:"not null"`
	Pay_Amount	  float64		`gorm:"type:decimal(16,2);not null"`
	Created_At	  time.Time
	Updated_At	  time.Time
}