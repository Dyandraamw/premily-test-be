package models

import (
	"fmt"
	"math/rand"
	_ "net/http"
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	Invoice_ID         string    `gorm:"size:100;uniqueIndex;not null;primary_key"`
	UserID             string    `gorm:"size:100"`
	Type               Type      `gorm:"not null"`
	Recipient          string    `gorm:"size:255;not null"`
	Address            string    `gorm:"type:text;not null"`
	Policy_Number      string    `gorm:"size:255;not null"`
	Name_Of_Insured    string    `gorm:"size:255;not null"`
	Address_Of_Insured string    `gorm:"size:255;not null"`
	Type_Of_Insurance  string    `gorm:"size:255;not null"`
	Period_Start       time.Time `gorm:"not null"`
	Period_End         time.Time `gorm:"not null"`
	Terms_Of_Period    string    `gorm:"size:255;not null"`
	Remarks            string    `gorm:"type:text;not null"`

	Description_Details Description_Details   `gorm:"foreignKey:Invoice_ID"`
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

	err = db.Debug().Model(&Invoice{}).Find(&invoice).Error
	if err != nil {

		return nil, err
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return &invoice, nil

}

func (i *Invoice) CreateInvoices(db *gorm.DB, invoices *Invoice) (*Invoice, error) {
	invoicesModels := &Invoice{
		Invoice_ID:         invoices.Invoice_ID,
		// UserID:             invoices.UserID,
		Type: invoices.Type,
		Recipient: invoices.Recipient,
		Address:            invoices.Address,
		Policy_Number:      invoices.Policy_Number,
		Name_Of_Insured:    invoices.Name_Of_Insured,
		Address_Of_Insured: invoices.Address_Of_Insured,
		Type_Of_Insurance:  invoices.Type_Of_Insurance,
		Period_Start:       invoices.Period_Start,
		Period_End:         invoices.Period_End,
		Terms_Of_Period:    invoices.Terms_Of_Period,
		Remarks:            invoices.Remarks,
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
