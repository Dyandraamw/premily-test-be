package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	_ "net/http"

	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	Invoice_ID           string    `gorm:"size:100;uniqueIndex;not null;primary_key"`
	UserID               string    `gorm:"size:100;default:''"`
	Company_Picture      string    `gorm:"size:255"`
	Company_Name         string    `gorm:"size:255;not null"`
	Company_Address      string    `gorm:"type:text;not null"`
	Company_Contact      string    `gorm:"size:255; not null"`
	Type                 Type      `gorm:"not null"`
	Recipient            string    `gorm:"size:255;not null"`
	Address              string    `gorm:"type:text;not null"`
	Net_Premium          Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Discount        Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Admin_Cost      Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Risk_Management Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_Brokage         Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Desc_PPH             Decimal   `gorm:"type:numeric(16,2);default:0.00"`
	Total_Premium_Due    Decimal   `gorm:"type:numeric(16,2);default:0.00"`
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

type ResponseInvoice struct {
	InvoiceID       string    `json:"Invoice_ID"`
	Recipient       string    `json:"Recipient"`
	CreatedAt       time.Time `json:"Created_At"`
	Period          string    `json:"Period"`
	TotalPremiumDue Decimal   `json:"Total_Premium_Due"`
}

// func calculateTotalDesc(invoice Invoice) decimal.Decimal {
//     descPremium, _ := decimal.NewFromString(invoice.Net_Premium.String())
//     descDiscount, _ := decimal.NewFromString(invoice.Desc_Discount.String())
//     descAdminCost, _ := decimal.NewFromString(invoice.Desc_Admin_Cost.String())
//     descRiskManagement, _ := decimal.NewFromString(invoice.Desc_Risk_Management.String())
//     descBrokage, _ := decimal.NewFromString(invoice.Desc_Brokage.String())
//     descPPH, _ := decimal.NewFromString(invoice.Desc_PPH.String())

//     total := descPremium.Add(descDiscount).
//              Add(descAdminCost).
//              Add(descRiskManagement).
//              Add(descBrokage).
//              Add(descPPH)
//     return total
// }

func (i *Invoice) GetInvoiceResponseList(db *gorm.DB) ([]ResponseInvoice, error) {
	var err error
	var invoices []Invoice

	err = db.Debug().Find(&invoices).Error
	if err != nil {
		return nil, err
	}

	responseInvoices := make([]ResponseInvoice, len(invoices))
	for idx, invoice := range invoices {
		period := invoice.Period_Start.Format("2006-01-02") + " - " + invoice.Period_End.Format("2006-01-02")
		responseInvoice := ResponseInvoice{
			InvoiceID:       invoice.Invoice_ID,
			Recipient:       invoice.Recipient,
			CreatedAt:       invoice.Created_At,
			Period:          period,
			TotalPremiumDue: Decimal{invoice.Total_Premium_Due.Decimal},
		}
		responseInvoices[idx] = responseInvoice
	}

	return responseInvoices, nil
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
	if err := db.Debug().First(&invoice, "Invoice_ID = ?", invoice_ID).Error; err != nil {
		return err
	}
	if err := db.Delete(&invoice).Error; err != nil {
		fmt.Printf("Fail!")
		return err
	}
	return nil
}

func (i *Invoice) UpdateInvoices(db *gorm.DB, invoiceID string, installments []Installment, sumInsured []Sum_Insured_Details) error {
	var invoice Invoice
	if err := db.First(&invoice, "invoice_id = ?", invoiceID).Error; err != nil {
		fmt.Println("invoice not found - model")
		return err
	}

	
	invoice.Recipient = i.Recipient
	invoice.Address = i.Address
	invoice.Net_Premium = i.Net_Premium
	invoice.Desc_Discount = i.Desc_Discount
	invoice.Desc_Admin_Cost = i.Desc_Admin_Cost
	invoice.Desc_Risk_Management = i.Desc_Risk_Management
	invoice.Desc_Brokage = i.Desc_Brokage
	invoice.Desc_PPH = i.Desc_PPH
	invoice.Total_Premium_Due = i.Total_Premium_Due
	invoice.Policy_Number = i.Policy_Number
	invoice.Name_Of_Insured = i.Name_Of_Insured
	invoice.Address_Of_Insured = i.Address_Of_Insured
	invoice.Type_Of_Insurance = i.Type_Of_Insurance
	invoice.Period_Start = i.Period_Start
	invoice.Period_End = i.Period_End
	invoice.Terms_Of_Period = i.Terms_Of_Period
	invoice.Remarks = i.Remarks
	invoice.Updated_At = i.Updated_At

	if err := db.Save(&invoice).Error; err != nil {
		return err
	}

	for _, inst := range installments {
		var existingInstallment Installment
		if err := db.Where("installment_id = ? AND invoice_id = ?", inst.Installment_ID, invoiceID).First(&existingInstallment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				inst.Invoice_ID = invoiceID
				if err := db.Create(&inst).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			existingInstallment.Due_Date = inst.Due_Date
			existingInstallment.Ins_Amount = inst.Ins_Amount
			if err := db.Save(&existingInstallment).Error; err != nil {
				return err
			}
		}
	}

	for _, sumIns := range sumInsured {
		var existingSumInsured Sum_Insured_Details
		if err := db.Where("sum_insured_id = ? AND invoice_id = ?", sumIns.Sum_Insured_ID, invoiceID).First(&existingSumInsured).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				sumIns.Invoice_ID = invoiceID
				if err := db.Create(&sumIns).Error; err != nil {
					return err
				}
			} else {
				existingSumInsured.Items_Name = sumIns.Items_Name
				existingSumInsured.Sum_Insured_Amount = sumIns.Sum_Insured_Amount
				existingSumInsured.Notes = sumIns.Notes
				if err := db.Save(&existingSumInsured).Error; err != nil {
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
		Company_Picture:      invoices.Company_Picture,
		Company_Name:         invoices.Company_Name,
		Company_Address:      invoices.Company_Address,
		Company_Contact:      invoices.Company_Contact,
		Type:                 invoices.Type,
		Recipient:            invoices.Recipient,
		Address:              invoices.Address,
		Net_Premium:          invoices.Net_Premium,
		Desc_Discount:        invoices.Desc_Discount,
		Desc_Admin_Cost:      invoices.Desc_Admin_Cost,
		Desc_Risk_Management: invoices.Desc_Risk_Management,
		Desc_Brokage:         invoices.Desc_Brokage,
		Desc_PPH:             invoices.Desc_PPH,
		Total_Premium_Due:    invoices.Total_Premium_Due,
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
