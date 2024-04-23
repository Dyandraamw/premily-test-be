package models

import "time"

type Statement_Of_Account struct {
	SOA_ID          string `gorm:"size:100;primary_key;not null"`
	User_ID         string `gorm:"size:100"`
	Name_Of_Insured string `gorm:"size:255;not null"`
	Period_Start   time.Time	`gorm:"not null"`
	Period_End		time.Time

	Statement_Of_Account_Details	[]Statement_Of_Account_Details	`gorm:"foreign_key:SOA_ID"`

	Created_At		time.Time
	Updated_At		time.Time
}