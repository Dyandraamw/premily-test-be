package models

import (
	"errors"
	"time"

	"log"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Adjustment struct {
	Adjustment_ID     string  `gorm:"size:100;primary_key;not null"`
	Payment_Status_ID string  `gorm:"size:100"`
	Adjustment_Title  string  `gorm:"size:255;not null"`
	Adjustment_Amount Decimal `gorm:"type:numeric(16,2);default:0;not null"`
	Created_At        time.Time
	Updated_At        time.Time
}

func (adjustment *Adjustment) CreateAdjustment(db *gorm.DB, adjust *Adjustment) (*Adjustment, error)  {
    adjustM := &Adjustment{
        Adjustment_ID: adjust.Adjustment_ID,
        Payment_Status_ID: adjust.Payment_Status_ID,
        Adjustment_Title: adjust.Adjustment_Title,
        Adjustment_Amount: adjust.Adjustment_Amount,
        Created_At: adjust.Created_At,
        Updated_At: adjust.Updated_At,
    }

    err := db.Debug().Create(&adjustM).Error
    if err != nil {
        return nil, err
    }
    return adjustM, nil
}

func GetTotalWithAdjustments(db *gorm.DB, selectedInvoiceID string, installmentID string, insAmount Decimal) (Decimal, error) {
    var invoice Invoice
    err := db.Preload("Installment").Where("invoice_id = ?", selectedInvoiceID).First(&invoice).Error
    if err != nil {
        log.Printf("Error loading invoice: %v", err)
        return Decimal{decimal.Zero}, err
    }

    // Cari installment yang sesuai dengan installmentID
    var targetInstallment Installment
    for _, installment := range invoice.Installment {
        if installment.Installment_ID == installmentID {
            targetInstallment = installment
            break
        }
    }

    // Jika installment tidak ditemukan
    if targetInstallment.Installment_ID == "" {
        log.Printf("Installment with ID %s not found in selected invoice %s", installmentID, selectedInvoiceID)
        return Decimal{decimal.Zero}, errors.New("Installment not found")
    }

    // Hitung total dengan penyesuaian
    total := targetInstallment.Ins_Amount
    var adjustments []Adjustment
    err = db.Where("installment_id = ?", installmentID).Find(&adjustments).Error
    if err != nil {
        log.Printf("Error loading adjustments: %v", err)
        return Decimal{decimal.Zero}, err
    }

    for _, adj := range adjustments {
        total = Decimal{total.Add(adj.Adjustment_Amount.Decimal)}
    }

    return total, nil
}
