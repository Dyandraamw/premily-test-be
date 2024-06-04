package models

import (
	"fmt"
	"time"

	"log"

	// "github.com/shopspring/decimal"
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

func CalculateAdjustment(db *gorm.DB, pStatID string) (Decimal, error)  {
    var adjustments []Adjustment

    // Get all adjustments for the given payment_status_id
    if err := db.Debug().Where("payment_status_id = ?", pStatID).Find(&adjustments).Error; err != nil {
        return Decimal{}, fmt.Errorf("error querying adjustments: %v", err)
    }

    // Debug: log jumlah adjustments yang diambil
    log.Printf("Number of adjustments found: %d\n", len(adjustments))

    // Calculate the total and row count
    rowCount := len(adjustments)
    var total Decimal

    if rowCount > 0 {
        for _, adj := range adjustments {
            log.Printf("Adding adjustment amount: %v\n", adj.Adjustment_Amount)
            total = Decimal{total.Add(adj.Adjustment_Amount.Decimal)}
        }
    } else {
        log.Printf("No adjustments found for payment_status_id: %s\n", pStatID)
    }

    return total, nil
}

func CalculatePayment(db *gorm.DB, pasStatID string)(Decimal, error){
    var payStat Payment_Status
    err := db.Debug().Preload("Invoice").Where("payment_status_id", pasStatID).First(&payStat).Error
    if err != nil {
        return Decimal{}, nil
    }
    invoiceID := payStat.Invoice_ID

    var installments []Installment
    err = db.Debug().Where("invoice_id", invoiceID).Find(&installments).Error
    if err != nil {
        return Decimal{}, nil
    }

    var payDet []Payment_Details

    for _, installment := range installments{
        var insPayDet []Payment_Details
        err := db.Debug().Where("installment_id", installment.Installment_ID).Find(&insPayDet).Error
        if err != nil {
            return Decimal{}, nil
        }
        payDet = append(payDet, insPayDet...)
    }

    rowCount := len(payDet)
    var balance Decimal
    if rowCount > 0{
        for _, pay := range payDet{
            balance = Decimal{balance.Add(pay.Pay_Amount.Decimal)}
            fmt.Println(balance)
        }
    }
    // balance = Decimal{balance.Sub(decimal.NewFromInt(int64(total)))}
    
    return balance, nil
}