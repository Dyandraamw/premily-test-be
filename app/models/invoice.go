package models

import (
	"fmt"
	"math/rand"
	_ "net/http"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Invoice struct {
	Invoice_ID           string          `gorm:"size:100;uniqueIndex;not null;primary_key"`
	UserID               string          `gorm:"size:100"`
	Type                 Type            `gorm:"not null"`
	Recipient            string          `gorm:"size:255;not null"`
	Address              string          `gorm:"type:text;not null"`
	Desc_Premium         decimal.Decimal `gorm:"type:numeric(16,2);"`
	Desc_Discount        decimal.Decimal `gorm:"type:numeric(16,2);"`
	Desc_Admin_Cost      decimal.Decimal `gorm:"type:numeric(16,2);"`
	Desc_Risk_Management decimal.Decimal `gorm:"type:numeric(16,2);"`
	Desc_Brokage         decimal.Decimal `gorm:"type:numeric(16,2);"`
	Desc_PPH             decimal.Decimal `gorm:"type:numeric(16,2);"`
	Policy_Number        string          `gorm:"size:255;not null"`
	Name_Of_Insured      string          `gorm:"size:255;not null"`
	Address_Of_Insured   string          `gorm:"size:255;not null"`
	Type_Of_Insurance    string          `gorm:"size:255;not null"`
	Period_Start         time.Time       `gorm:"not null"`
	Period_End           time.Time       `gorm:"not null"`
	Terms_Of_Period      string          `gorm:"size:255;not null"`
	Remarks              string          `gorm:"type:text;not null"`

	Installment         []Installment         `gorm:"foreignKey:Invoice_ID"`
	Sum_Insured_Details []Sum_Insured_Details `gorm:"foreignKey:Invoice_ID"`
	Payment_Status      Payment_Status        `gorm:"foreignKey:Invoice_ID"`

	Created_At time.Time
	Updated_At time.Time
}

type Type string

const (
	CreditType = "credit"
	DebitType  = "debit"
)

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

func (i *Invoice) UpdateInvoices() {

}

func (i *Invoice) DeletedInvoices(db *gorm.DB, invoice_ID string) error{
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