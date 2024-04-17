package models

import "time"

type Invoice struct {
	InvoiceID          string `gorm:"size:100;uniqueIndex;not null;primary_key"`
	User               User
	UserID             string `gorm:"size:100;index"`
	Type               Type   `gorm:"type:enum('credit','debit'); not null"`
	Recipient          string `gorm:"size:255;not null"`
	Address            string `gorm:"type:text"`
	Policy_Number      string `gorm:"size:255;not null"`
	Name_Of_Insured    string `gorm:"size:255;not null"`
	Address_Of_Insured string `gorm:"size:255;not null"`
	Type_Of_Insurance  string `gorm:"size:255;not null"`
	Periode_Start      time.Time
	Periode_End			time.Time
	Terms_Of_Periode	string	`gorm:"size:255;not null"`
	Remarks				string	`gorm:"type:text;not null"`
}

type Type string

const (
	CreditType = "credit"
	DebitType  = "debit"
)