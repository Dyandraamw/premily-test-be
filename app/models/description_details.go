package models

type Description_Details struct{
	Desc_Details_ID		string	`gorm:"size:100;not null; primary_key"`
	Invoice_ID			string	`gorm:"size:100;"`
	Premium				float64	`gorm:"type:decimal(16,2);not null"`
	Discount			float64	`gorm:"type:decimal(16,2);not null"`
	Admin_Cost			float64	`gorm:"type:decimal(16,2;not null)"`
	Risk_Management		float64	`gorm:"type:decimal(16,2);"`
	Brokage				float64	`gorm:"type:decimal(16,2);"`
	PPH					float64	`gorm:"type:decimal(16,2);"`
}