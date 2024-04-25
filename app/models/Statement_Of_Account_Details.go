package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Statement_Of_Account_Details struct {
	SOA_Details_ID      	 string `gorm:"size:100;primary_key;not null"`
	SOA_ID              	 string `gorm:"size:100"`
	Invoice_ID          	 string `gorm:"size:100"`
	Recipient            	string `gorm:"size:255;not null"`
	Installment_Standing 	uint   `gorm:"not null"`
	Due_Date             	time.Time		`gorm:"not null"`
	SOA_Amount				decimal.Decimal		`gorm:"type:numeric(16,2);not null"`
	Payment_Date		 	time.Time	`gorm:"not null"`
	Payment_Amount		 decimal.Decimal		`gorm:"type:numeric(16,2);not null"`
	Payment_Allocation	 decimal.Decimal		`gorm:"type:numeric(16,2);not null"`
	Status				 string		`gorm:"size:255;not null"`
	Aging				 uint		`gorm:"not null"`
	Created_At			 time.Time
	Updated_At			 time.Time
}