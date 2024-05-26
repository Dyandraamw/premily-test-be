package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	_ "github.com/gorilla/mux"
)

func (server *Server) CreateNewPaymentStatus(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// payment_S_ID := vars["payment_status_id"]

	invoiceModel := models.Invoice{}
	_, err := invoiceModel.GetInvoice(server.DB)
	if err != nil {
		http.Error(w, "Retrive Invoice fail - payment", http.StatusBadRequest)
		return
	}

	selectedInvoiceID := r.FormValue("invoice_id")
	if selectedInvoiceID == "" {
		http.Error(w, "Please select invoice!", http.StatusBadRequest)
		return
	}

	var idGeneratorPaymentS = NewIDGenerator("Payment")

	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	paymentS_M := models.Payment_Status{}

	payment := &models.Payment_Status{}
	for {
		paymentID := idGeneratorPaymentS.NextID()

		var existPayments models.Payment_Status
		scan := server.DB.Where("payment_status_id = ? ", paymentID).First(&existPayments)
		if scan.RowsAffected == 0 {
			payment = &models.Payment_Status{
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

	// var paymentShow models.Payment_Status
	payShow, err := payment.CreateNewPayment(server.DB, &paymentS_M)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mengambil data Installment terkait dengan invoice
	// installments := models.Installment{}
	// installmentList, err := installments.GetInstallmentByInvoiceID(server.DB, (*invoices)[0].Invoice_ID)
	// if err != nil {
	// 	http.Error(w, "Retrive fails", http.StatusBadRequest)
	// 	return
	// }

	// data, err := json.Marshal(installmentList)
	// if err != nil {
	// 	http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
	// 	return
	// }
	data, _ := json.Marshal(payShow)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	// Set response headers and write JSON data
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Successfully!"))

	// Membentuk response JSON
	// response := struct {
	// 	PaymentStatus models.Payment_Status `json:"payment_status"`
	// 	Installments  []models.Installment  `json:"installments"`
	// }{
	// 	PaymentStatus: paymentShow,
	// 	Installments:  *installmentList,
	// }

	// // Mengirim response JSON ke client
	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(response); err != nil {
	// 	http.Error(w, "Encode response fail", http.StatusInternalServerError)
	// 	return
	// }

}
