package models

import (
	"fmt"
	"time"

	// "github.com/go-logr/logr/funcr"
	"gorm.io/gorm"
)

type Statement_Of_Account_Details struct {
	SOA_Details_ID       string    `gorm:"size:100;primary_key;not null"`
	SOA_ID               string    `gorm:"size:100"`
	Invoice_ID           string    `gorm:"size:100"`
	Recipient            string    `gorm:"size:255;not null;default:''"`
	Installment_Standing uint      `gorm:"not null"`
	Due_Date             time.Time `gorm:"not null;default:current_timestamp"`
	SOA_Amount           Decimal   `gorm:"type:numeric(16,2);default:0;not null"`
	Payment_Date         time.Time `gorm:"not null"`
	Payment_Amount       Decimal   `gorm:"type:numeric(16,2);default:0;not null"`
	Payment_Allocation   Decimal   `gorm:"type:numeric(16,2);default:0;not null"`
	Status               string    `gorm:"size:255;not null"`
	Aging                uint      `gorm:"not null"`
	Created_At           time.Time
	Updated_At           time.Time
}

func (soa_d *Statement_Of_Account_Details) CreateSoaDetails(db *gorm.DB, soa_details *Statement_Of_Account_Details) (*Statement_Of_Account_Details, error) {
	soa_Details_Model := &Statement_Of_Account_Details{
		SOA_Details_ID:       soa_details.SOA_Details_ID,
		SOA_ID:               soa_details.SOA_ID,
		Invoice_ID:           soa_details.Invoice_ID,
		Recipient:            soa_details.Recipient,
		Installment_Standing: soa_details.Installment_Standing,
		Due_Date:             soa_details.Due_Date,
		SOA_Amount:           soa_details.SOA_Amount,
		Payment_Date:         soa_details.Payment_Date,
		Payment_Amount:       soa_details.Payment_Amount,
		Payment_Allocation:   soa_details.Payment_Allocation,
		Status:               soa_details.Status,
		Aging:                soa_details.Aging,
		Created_At:           soa_details.Created_At,
		Updated_At:           soa_details.Updated_At,
	}

	err := db.Debug().Create(&soa_Details_Model).Error
	if err != nil {
		return nil, err
	}
	return soa_Details_Model, nil
}

func (soa_d *Statement_Of_Account_Details) UpdatesItemsSoa(db *gorm.DB, soa_id string) error {
	if soa_d == nil {
		return fmt.Errorf("received nil Statement_Of_Account_Details")
	}

	var items Statement_Of_Account_Details

	if err := db.First(&items, "soa_details_id = ?", soa_id).Error; err != nil {
		return fmt.Errorf("items not found: %w", err)
	}

	items.Recipient = soa_d.Recipient
	items.Installment_Standing = soa_d.Installment_Standing
	items.Due_Date = soa_d.Due_Date
	items.SOA_Amount = soa_d.SOA_Amount
	items.Payment_Date = soa_d.Payment_Date
	items.Payment_Amount = soa_d.Payment_Amount
	items.Updated_At = soa_d.Updated_At

	if err := db.Save(&items).Error; err != nil {
		return fmt.Errorf("failed to save updated items: %w", err)
	}

	return nil
}


func (item *Statement_Of_Account_Details) GetItemsBySoaID(db *gorm.DB, soa_id string)(*[]Statement_Of_Account_Details, error){
	var items []Statement_Of_Account_Details
	err := db.Debug().Where("soa_id = ?", soa_id).Find(&items).Error
	if err != nil {
		fmt.Println("Retrive items fail - model"+err.Error())
		return nil, err
	}

	return &items, nil
}