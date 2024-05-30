package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	_ "github.com/gorilla/mux"
	
)

func (server *Server) CreatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	invoiceModel := models.Invoice{}
	invoices, err := invoiceModel.GetInvoice(server.DB)
	if err != nil {
		http.Error(w, "Retrive Invoice fail - payment", http.StatusBadRequest)
		return
	}

	// Display all invoices
	if r.Method == http.MethodGet {
		response, _ := json.Marshal(invoices)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

	selectedInvoiceID := r.FormValue("invoice_id")
	if selectedInvoiceID == "" {
		http.Error(w, "Please select an invoice!", http.StatusBadRequest)
		return
	}

	log.Printf("Selected Invoice ID: %s", selectedInvoiceID)

	// Cari invoice yang dipilih
	var selectedInvoice *models.Invoice
	for x, invoice := range *invoices {
		if invoice.Invoice_ID == selectedInvoiceID {
			selectedInvoice = &(*invoices)[x]
			break
		}
	}

	// Jika invoice yang dipilih tidak ditemukan
	if selectedInvoice == nil {
		http.Error(w, "Selected invoice not found!", http.StatusBadRequest)
		return
	}

	_, err = invoiceModel.GetInvoiceByIDmodel(server.DB, selectedInvoiceID)
	if err != nil {
		http.Error(w, "Invoice not found!", http.StatusNotFound)
		return
	}

	var idGeneratorPaymentS = NewIDGenerator("Payment")

	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paymentS := &models.Payment_Status{}
	for {
		paymentID := idGeneratorPaymentS.NextID()

		var existPayments models.Payment_Status
		scan := server.DB.Where("payment_status_id = ? ", paymentID).First(&existPayments)
		if scan.RowsAffected == 0 {
			paymentS = &models.Payment_Status{
				Payment_Status_ID: paymentID,
				UserID:            userID,
				Invoice_ID:        selectedInvoiceID,
				Status:            "PENDING",
				Created_At:        time.Now(),
				Updated_At:        time.Now(),
			}
			break
		}
	}

	// payment_M := models.Payment_Status{}
	_, err = paymentS.CreateNewPayment(server.DB)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Add Payment fail-control", http.StatusBadRequest)
		return
	}

	// Query invoice dari database berdasarkan selectedInvoiceID
	var invoice models.Invoice
	err = server.DB.Preload("Installment").First(&invoice, "invoice_id = ?", selectedInvoiceID).Error
	if err != nil {
		http.Error(w, "Error fetching invoice data", http.StatusInternalServerError)
		return
	}

	// Menyiapkan data response untuk installments
	type InstallmentInfo struct {
		InstallmentID string         `json:"installment_id"`
		DueDate       time.Time      `json:"due_date"`
		InsAmount     models.Decimal `json:"ins_amount"`
	}

	var installmentInfos []InstallmentInfo
	for _, installment := range invoice.Installment {
		installmentInfos = append(installmentInfos, InstallmentInfo{
			InstallmentID: installment.Installment_ID,
			DueDate:       installment.Due_Date,
			InsAmount:     installment.Ins_Amount,
		})
	}

	// Marshal data response ke JSON
	responseData := struct {
		InvoiceID       string            `json:"invoice_id"`
		Recipient       string            `json:"recipient"`
		Period_Start    time.Time         `json:"periode_start"`
		Periode_End     time.Time         `json:"periode_end"`
		InstallmentInfo []InstallmentInfo `json:"installment_info"`
	}{
		InvoiceID:       selectedInvoiceID,
		Recipient:       invoice.Recipient,
		Period_Start:    invoice.Period_Start,
		Periode_End:     invoice.Period_End,
		InstallmentInfo: installmentInfos,
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}

	// Set header dan kirim response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

