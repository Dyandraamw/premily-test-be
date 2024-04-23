package models

import "time"

type Payment_Status struct {
	Payment_Status_ID string `gorm:"size:100;primary_key;not null"`
	UserID            string `gorm:"size:100"`
	Invoice_ID        string `gorm:"size:100"`
	Status            string `gorm:"size:255;not null"`

	Adjustment			[]Adjustment	`gorm:"foreign_key:Payment_Status_ID"`

	Created_At        time.Time
	Updated_At		time.Time
}