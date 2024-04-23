package models

import "time"

type Invoice struct {
	Invoice_ID          string `gorm:"size:100;uniqueIndex;not null;primary_key"`
	UserID				string	`gorm:"size:100"`
	Type               Type   `gorm:"not null"`
	Recipient          string `gorm:"size:255;not null"`
	Address            string `gorm:"type:text;not null"`
	Policy_Number      string `gorm:"size:255;not null"`
	Name_Of_Insured    string `gorm:"size:255;not null"`
	Address_Of_Insured string `gorm:"size:255;not null"`
	Type_Of_Insurance  string `gorm:"size:255;not null"`
	Period_Start       time.Time	`gorm:"not null"`
	Period_End		   time.Time	`gorm:"not null"`
	Terms_Of_Period		string	`gorm:"size:255;not null"`
	Remarks				string	`gorm:"type:text;not null"`

	Description_Details	Description_Details		`gorm:"foreign_key:Invoice_ID"`
	Installment			[]Installment			`gorm:"foreign_key:Invoice_ID"`
	Sum_Insured_Details	[]Sum_Insured_Details	`gorm:"foreign_key:Invoice_ID"`
	Payment_Status		Payment_Status			`gorm:"foreign_key:Invoice_ID"`

	Created_At		time.Time
	Updated_At		time.Time
}

type Type string

const (
	CreditType = "credit"
	DebitType  = "debit"
)