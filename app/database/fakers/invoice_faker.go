package fakers

import (
	"fmt"
	"time"

	_ "github.com/bxcodec/faker/v4"
	"github.com/frangklynndruru/premily_backend/app/models"
	"gorm.io/gorm"
)

var idGenerator = NewIDGenerator("D")

func InvoiceFaker(db *gorm.DB) *models.Invoice {
	for {
		// Dapatkan ID baru dari generator
		invoiceID := idGenerator.NextID()

		// Cek apakah ID sudah ada di dalam database
		var existingInvoice models.Invoice
		result := db.Where("invoice_id = ?", invoiceID).First(&existingInvoice)
		if result.RowsAffected == 0 {
			// Jika ID belum ada di dalam database, buat objek Invoice dan simpan
			newInvoice := &models.Invoice{
				Invoice_ID:         invoiceID,
				UserID:            "ebde4360-92d4-4c3b-ac36-4a065ed912a6",
				Type:               "Debit",
				Recipient:          "PT Garuda Indonesia",
				Address:            "",
				Policy_Number:      "",
				Name_Of_Insured:    "PT Gunung Barat",
				Address_Of_Insured: "",
				Type_Of_Insurance:  "Pendidikan",
				Period_Start:       time.Now(),
				Period_End:         time.Now(),
				Terms_Of_Period:    "Based on AVN 6A- the premium shall be paid in the following instalments :",

				Created_At: time.Now(),
				Updated_At: time.Now(),
			}

			// Simpan objek Invoice ke dalam database
			db.Create(newInvoice)

			// Kembalikan objek Invoice yang baru dibuat
			return newInvoice
		}
	}
}

type IDGenerator struct {
	prefix string
	count  int
}

func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{
		prefix: prefix,
		count:  0,
	}
}

func (g *IDGenerator) NextID() string {
	g.count = g.count + 1
	return fmt.Sprintf("%s-%s", g.prefix, g.pad(g.count, 5))
}

func (g *IDGenerator) pad(number, width int) string {
	return fmt.Sprintf("%0*d", width, number)
}
