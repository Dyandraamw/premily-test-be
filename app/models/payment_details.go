package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment_Details struct {
	Pay_Detail_ID 	string 		`gorm:"size:100;primary_key;not null"`
	Installment_ID	string		`gorm:"size:100"`
	Pay_Date      	time.Time		`gorm:"not null;default:current_timestamp"`
	Pay_Amount	  	Decimal		`gorm:"type:numeric(16,2);default:0;not null"`
	Created_At	  	time.Time
	Updated_At	  	time.Time
}

func (pDtl *Payment_Details) CreatePaymentDetails(db *gorm.DB, pDetails *Payment_Details)(*Payment_Details, error){
	payment_detailsM := &Payment_Details{
		Pay_Detail_ID: pDetails.Pay_Detail_ID,
		Installment_ID: pDetails.Installment_ID,
		Pay_Date: pDetails.Pay_Date,
		Pay_Amount: pDetails.Pay_Amount,
		Created_At: pDetails.Created_At,
		Updated_At: pDetails.Updated_At,
	}

	err := db.Debug().Create(&payment_detailsM).Error
	if err != nil{
		return nil, err
	}
	return payment_detailsM, nil
}

