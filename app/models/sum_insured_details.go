package models

type Sum_Insured_Details struct{
	Sum_Insured_ID		string	`gorm:"size:100;not null;primary_key"`
	Invoice_ID			string	`gorm:"size:100"`
	Items_Name				string	`gorm:"size:255;not null"`
	Sum_Insured_Amount	float64	`gorm:"type:decimal(16,2);not null"`
	Notes				string	`gorm:"size:255;not null"`
}