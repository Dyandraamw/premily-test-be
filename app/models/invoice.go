package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	_ "net/http"

	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Invoice struct {
	Invoice_ID           string    `gorm:"size:100;uniqueIndex;not null;primary_key"`
	UserID               string    `gorm:"size:100;default:''"`
	Type                 Type      `gorm:"not null"`
	Recipient            string    `gorm:"size:255;not null"`
	Address              string    `gorm:"type:text;not null"`
	Desc_Premium         Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Discount        Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Admin_Cost      Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Risk_Management Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Brokage         Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_PPH             Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Policy_Number        string    `gorm:"size:255;not null;default:''"`
	Name_Of_Insured      string    `gorm:"size:255;not null;default:''"`
	Address_Of_Insured   string    `gorm:"size:255;not null;default:''"`
	Type_Of_Insurance    string    `gorm:"size:255;not null;default:''"`
	Period_Start         time.Time `gorm:"not null;default:current_timestamp"`
	Period_End           time.Time `gorm:"not null;default:current_timestamp"`
	Terms_Of_Period      string    `gorm:"size:255;not null;default:''"`
	Remarks              string    `gorm:"type:text;not null;default:''"`

	Installment         []Installment         `gorm:"foreignKey:Invoice_ID;constraint:OnDelete:CASCADE"`
	Sum_Insured_Details []Sum_Insured_Details `gorm:"foreignKey:Invoice_ID;constraint:OnDelete:CASCADE"`
	Payment_Status      Payment_Status        `gorm:"foreignKey:Invoice_ID"`

	Created_At time.Time
	Updated_At time.Time
}

type Type string

const (
	CreditType = "credit"
	DebitType  = "debit"
)

func calculateTotalDesc(invoice Invoice) decimal.Decimal {
    descPremium, _ := decimal.NewFromString(invoice.Desc_Premium.String())
    descDiscount, _ := decimal.NewFromString(invoice.Desc_Discount.String())
    descAdminCost, _ := decimal.NewFromString(invoice.Desc_Admin_Cost.String())
    descRiskManagement, _ := decimal.NewFromString(invoice.Desc_Risk_Management.String())
    descBrokage, _ := decimal.NewFromString(invoice.Desc_Brokage.String())
    descPPH, _ := decimal.NewFromString(invoice.Desc_PPH.String())

    total := descPremium.Add(descDiscount).
             Add(descAdminCost).
             Add(descRiskManagement).
             Add(descBrokage).
             Add(descPPH)
    return total
}

func (i *Invoice) GetInvoice(db *gorm.DB) (*[]Invoice, error) {
	var err error
	var invoice []Invoice

	err = db.Debug().Preload("Installment").Preload("Sum_Insured_Details").Find(&invoice).Error
	if err != nil {

		return nil, err
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return &invoice, nil

}
func (i *Invoice) GetInvoiceByIDmodel(db *gorm.DB, invoice_ID string) (*Invoice, error) {
	var err error
	var invoice Invoice

	err = db.Debug().Model(&Invoice{}).Preload("Installment").Preload("Sum_Insured_Details").First(&invoice, "invoice_id = ?", invoice_ID).Error
	if err != nil {
		log.Printf("Error fetching invoice by ID: %v", err)
		return nil, err
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return &invoice, nil

}

func (i *Invoice) DeletedInvoices(db *gorm.DB, invoice_ID string) error {
	invoice := &Invoice{}
	if err := db.Debug().First(&invoice, "invoice_id = ?", invoice_ID).Error; err != nil {
		return err
	}
	if err := db.Delete(&invoice).Error; err != nil {
		fmt.Printf("Fail!")
		return err
	}
	return nil
}

func (i *Invoice) UpdateInvoices(db *gorm.DB, invoice_ID string, installments []Installment, sum_insured []Sum_Insured_Details) error {
	var invoice Invoice
	if err := db.First(&invoice, "invoice_id = ?", invoice_ID).Error; err != nil {
		fmt.Println("invoice tidak ditemukan - model")
		return nil
	}
	invoice.Type = i.Type
	invoice.Recipient = i.Recipient
	invoice.Address = i.Address
	invoice.Desc_Premium = i.Desc_Premium
	invoice.Desc_Discount = i.Desc_Discount
	invoice.Desc_Admin_Cost = i.Desc_Admin_Cost
	invoice.Desc_Risk_Management = i.Desc_Risk_Management
	invoice.Desc_Brokage = i.Desc_Brokage
	invoice.Desc_PPH = i.Desc_PPH
	invoice.Policy_Number = i.Policy_Number
	invoice.Name_Of_Insured = i.Name_Of_Insured
	invoice.Address_Of_Insured = i.Address_Of_Insured
	invoice.Type_Of_Insurance = i.Type_Of_Insurance
	invoice.Period_Start = i.Period_Start
	invoice.Period_End = i.Period_End
	invoice.Terms_Of_Period = i.Terms_Of_Period
	invoice.Remarks = i.Remarks
	invoice.Created_At = i.Created_At
	invoice.Updated_At = i.Updated_At

	err := db.Save(&invoice).Error
	if err != nil {
		return err
	}

	for _, installment := range installments {
		var existInstallment Installment
		err := db.Where("installment_id = ? AND invoice_id = ?", installment.Installment_ID, invoice_ID).First(&existInstallment).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				installment.Invoice_ID = invoice_ID
				if err := db.Create(&installment).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			existInstallment.Due_Date = installment.Due_Date
			existInstallment.Ins_Amount = installment.Ins_Amount
			if err := db.Save(&installment).Error; err != nil {
				return err
			}
		}
	}

	for _, sum_ins := range sum_insured {
		var existSumIns Sum_Insured_Details
		err := db.Where("Sum_Insured_ID = ? AND invoice_id", sum_ins.Sum_Insured_ID, invoice_ID).First(&existSumIns).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				sum_ins.Invoice_ID = invoice_ID
				if err := db.Create(&sum_ins).Error; err != nil {
					return err
				} else {
					return err
				}
			} else {
				existSumIns.Items_Name = sum_ins.Items_Name
				existSumIns.Sum_Insured_Amount = sum_ins.Sum_Insured_Amount
				existSumIns.Notes = sum_ins.Notes
				err := db.Save(&sum_ins).Error
				if err != nil {
					return err
				}
			}
		}
	}

	return nil

}

func (i *Invoice) CreateInvoices(db *gorm.DB, invoices *Invoice) (*Invoice, error) {
	invoicesModels := &Invoice{
		Invoice_ID:           invoices.Invoice_ID,
		UserID:               invoices.UserID,
		Type:                 invoices.Type,
		Recipient:            invoices.Recipient,
		Address:              invoices.Address,
		Desc_Premium:         invoices.Desc_Premium,
		Desc_Discount:        invoices.Desc_Discount,
		Desc_Admin_Cost:      invoices.Desc_Admin_Cost,
		Desc_Risk_Management: invoices.Desc_Risk_Management,
		Desc_Brokage:         invoices.Desc_Brokage,
		Desc_PPH:             invoices.Desc_PPH,
		Policy_Number:        invoices.Policy_Number,
		Name_Of_Insured:      invoices.Name_Of_Insured,
		Address_Of_Insured:   invoices.Address_Of_Insured,
		Type_Of_Insurance:    invoices.Type_Of_Insurance,
		Period_Start:         invoices.Period_Start,
		Period_End:           invoices.Period_End,
		Terms_Of_Period:      invoices.Terms_Of_Period,
		Remarks:              invoices.Remarks,
		Created_At:           invoices.Created_At,
		Updated_At:           invoices.Updated_At,
	}
	err := db.Debug().Create(&invoicesModels).Error
	if err != nil {
		return nil, err
	}

	return invoicesModels, nil
}

func GenerateInvoiceID(db *gorm.DB, type_of_invoices Type) (string, error) {

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)
	prefix := ""

	if type_of_invoices == "debit" {
		prefix = "DN"
	} else if type_of_invoices == "credit" {
		prefix = "CN"
	} else {
		prefix = "UNK"
	}

	return fmt.Sprintf("%s-%05d", prefix, randomNumber), nil
}


