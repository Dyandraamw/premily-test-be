package models

import "time"

type Adjustment struct {
	Adjustment_ID     string  `gorm:"size:100;primary_key;not null"`
	Payment_Status_ID string  `gorm:"size:100"`
	Adjustment_Title  string  `gorm:"size:255;not null"`
	Adjustment_Amount float64 `gorm:"type:decimal(16,2);not null"`
	Created_At        time.Time
	Updated_At			time.Time
}