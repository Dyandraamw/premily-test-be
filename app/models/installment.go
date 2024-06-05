package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Installment struct {
	Installment_ID string    `gorm:"size:100;not null;primary_key"`
	Invoice_ID     string    `gorm:"size:100"`
	Due_Date       time.Time `gorm:"not null;default:current_timestamp"`
	Ins_Amount     Decimal   `gorm:"type:numeric(16,2);default:0;not null"`

	Payment_Details []Payment_Details `gorm:"foreignKey:Installment_ID;constraint:OnDelete:CASCADE"`
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

func (in *Installment) GetInstallmentByInvoiceID(db *gorm.DB, invoiceID string) (*[]Installment, error) {
	var installments []Installment
	err := db.Debug().Where("invoice_id = ?", in.Invoice_ID).Find(&installments).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &installments, nil
}

func calculateTotalInstallment(db *gorm.DB, invoiceID string) Decimal {
    var totalInstallment Decimal
 
    var sumInstallment struct {
        TotalInstallment Decimal
    }
    if err := db.Raw("SELECT SUM(ins_amount) AS total_installment FROM installments WHERE invoice_id = ?", invoiceID).Scan(&sumInstallment).Error; err != nil {
        log.Fatalf("Failed to calculate total installment: %v", err)
        
		return Decimal{}
    }
    totalInstallment = sumInstallment.TotalInstallment
    return totalInstallment
}

