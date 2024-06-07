package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Payment_Details struct {
	Pay_Detail_ID  string    `gorm:"size:100;primary_key;not null"`
	Installment_ID string    `gorm:"size:100"`
	Pay_Date       time.Time `gorm:"not null;default:current_timestamp"`
	Pay_Amount     Decimal   `gorm:"type:numeric(16,2);default:0;not null"`
	Created_At     time.Time
	Updated_At     time.Time
}

func (pDtl *Payment_Details) CreatePaymentDetails(db *gorm.DB, pDetails *Payment_Details) (*Payment_Details, error) {
	payment_detailsM := &Payment_Details{
		Pay_Detail_ID:  pDetails.Pay_Detail_ID,
		Installment_ID: pDetails.Installment_ID,
		Pay_Date:       pDetails.Pay_Date,
		Pay_Amount:     pDetails.Pay_Amount,
		Created_At:     pDetails.Created_At,
		Updated_At:     pDetails.Updated_At,
	}

	err := db.Debug().Create(&payment_detailsM).Error
	if err != nil {
		return nil, err
	}
	return payment_detailsM, nil
}

func (pay *Payment_Details) UpdatePayment(db *gorm.DB, pay_det_id string) error {
	var payDet Payment_Details

	err := db.Debug().First(&payDet, "pay_detail_id = ?", pay_det_id).Error
	if err != nil {
		return fmt.Errorf("Details of payment not found ! %w (model)", err)
	}
	payDet.Pay_Date = pay.Pay_Date
	payDet.Pay_Amount = pay.Pay_Amount
	payDet.Updated_At = pay.Updated_At

	if err := db.Save(&payDet).Error; err != nil{
		return fmt.Errorf("failed to save updated detail of payment: %w", err)
	}

	return nil
}
