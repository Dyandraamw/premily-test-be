package models

import (
	"gorm.io/gorm"
	"fmt"
)

type Sum_Insured_Details struct {
	Sum_Insured_ID     string  `gorm:"size:100;not null;primary_key"`
	Invoice_ID         string  `gorm:"size:100"`
	Items_Name         string  `gorm:"size:255;not null"`
	Sum_Insured_Amount Decimal `gorm:"type:numeric(16,2);default:0;not null"`
	Notes              string  `gorm:"size:255;not null"`
}

func (SumIns *Sum_Insured_Details) CreateSumInsuredDetails(db *gorm.DB, sumInsured *Sum_Insured_Details) (*Sum_Insured_Details, error) {
	sumInsuredModels := &Sum_Insured_Details{
		Sum_Insured_ID:     sumInsured.Sum_Insured_ID,
		Invoice_ID:         sumInsured.Invoice_ID,
		Items_Name:         sumInsured.Items_Name,
		Sum_Insured_Amount: sumInsured.Sum_Insured_Amount,
		Notes:              sumInsured.Notes,
	}

	err := db.Debug().Create(&sumInsuredModels).Error
	if err != nil {
		return nil, err
	}
	return sumInsuredModels, nil
}

func (in *Installment) GetSumInsByInvoiceID(db *gorm.DB, invoiceID string) (*[]Sum_Insured_Details, error) {
	var sum_ins []Sum_Insured_Details
	err := db.Debug().Where("invoice_id = ?", invoiceID).Find(&sum_ins).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &sum_ins, nil
}

