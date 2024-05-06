package fakers

import (
	"fmt"
	"time"

	_"github.com/bxcodec/faker/v4"
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
                UserID:             "98d92fbb-85fa-4f13-bd9a-d6e011332b56",
                Type:               "Debit",
                Recipient:          "PT Garuda Indonesia",
                Address:            "",
                Policy_Number:      "",
                Name_Of_Insured:    "PT Gunung Barat",
                Address_Of_Insured: "",
                Type_Of_Insurance:  "Pendidikan",
                Period_Start:       time.Time{},
                Period_End:         time.Time{},
                Terms_Of_Period:    "Based on AVN 6A- the premium shall be paid in the following instalments :",

                Created_At: time.Time{},
                Updated_At: time.Time{},
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
	return fmt.Sprintf("%s-%s", g.prefix, g.pad(g.count, 3))
}

func (g *IDGenerator) pad(number, width int) string{
	return fmt.Sprintf("%0*d", width, number)
}

