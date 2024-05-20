package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Installment struct {
	Installment_ID string          `gorm:"size:100;not null;primary_key"`
	Invoice_ID     string          `gorm:"size:100"`
	Due_Date       time.Time       `gorm:"not null"`
	Ins_Amount     decimal.Decimal `gorm:"type:numeric(16,2);not null"`

	Payment_Details []Payment_Details `gorm:"foreignKey:Installment_ID"`
}

func (ins *Installment) CreateInstallment(db *gorm.DB, installment *Installment) (*Installment, error) {
	installmentModels := &Installment{
		Installment_ID: installment.Installment_ID,
		Invoice_ID:     installment.Invoice_ID,
		Due_Date:       installment.Due_Date,
		Ins_Amount:     installment.Ins_Amount,
	}

	err := db.Debug().Create(&installmentModels).Error
	if err != nil {
		return nil, err
	}
	return installmentModels, nil
}

func (in *Installment) GetInstallmentByInvoiceID(db *gorm.DB, invoiceID string) (*[]Installment, error){
	var installments []Installment
	err := db.Debug().Where("invoice_id = ?", invoiceID).Find(&installments).Error
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}
	return &installments, nil
}