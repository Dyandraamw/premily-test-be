package models

import (
	"time"

	
	"gorm.io/gorm"
)

type Payment_Status struct {
	Payment_Status_ID string `gorm:"size:100;primary_key;not null"`
	UserID            string `gorm:"size:100"`
	Invoice_ID        string `gorm:"size:100"`
	Status            string `gorm:"size:255;not null"`

	Adjustment []Adjustment `gorm:"foreignKey:Payment_Status_ID;constrain:OnUpdate, OnDelete:CASCADE"`

	Created_At time.Time
	Updated_At time.Time
}

func (p *Payment_Status) CreateNewPayment(db *gorm.DB, paymentS *Payment_Status) (*Payment_Status, error){
	payment_Model := &Payment_Status{
		Payment_Status_ID: paymentS.Payment_Status_ID,
		UserID: paymentS.UserID,
		Invoice_ID: paymentS.Invoice_ID,
		Status: paymentS.Status,
		Created_At: paymentS.Created_At,
		Updated_At: paymentS.Updated_At,
	}
	invoice := &Invoice{}
	err := db.Debug().Preload("installment").First(&invoice, "invoice_id = ? AND type=?", invoice.Invoice_ID, invoice).Error
	if err != nil{
		return nil, err
	}

	err = db.Debug().Create(&payment_Model).Error
	if err != nil {
		return nil, err
	}
	return payment_Model, nil
}
