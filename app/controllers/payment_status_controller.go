package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	_ "github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

/*
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
*/

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
		 InstallmentInfo []InstallmentInfo `json:"installment_info"`
	 }{
		 InvoiceID:       selectedInvoiceID,
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



func (server *Server) AddPayment(w http.ResponseWriter, r *http.Request) {
	installmentID := r.FormValue("installment_id")
	p_date_detail := r.FormValue("pay_date")
	p_amount_detail := r.FormValue("pay_amount")

	if installmentID == "" || p_date_detail == "" || p_amount_detail == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}

	var installment models.Installment
	err := server.DB.Preload("Payment_Details").First(&installment, "installment_id = ?", installmentID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (decimal.Decimal, error) {
		return decimal.NewFromString(s)
	}

	const layoutTime = "02-01-2006"
	pay_date, err := time.Parse(layoutTime, p_date_detail)
	if err != nil {
		http.Error(w, "invalid payment date!", http.StatusBadRequest)
		return
	}

	pay_amount, err := convertToDecimal(p_amount_detail)
	if err != nil {
		http.Error(w, "invalid payment amount!", http.StatusBadRequest)
		return
	}

	var prevBalance decimal.Decimal
	for _, detail := range installment.Payment_Details {
		prevBalance = prevBalance.Add(detail.Pay_Amount.Decimal)
	}

	var paymentAllocation decimal.Decimal
	if prevBalance.LessThan(decimal.Zero) {
		paymentAllocation = pay_amount
	} else {
		paymentAllocation = prevBalance.Sub(pay_amount)
	}

	var idGeneratorPaymentDetail = NewIDGenerator("Payment-Details")

	pay_detail := &models.Payment_Details{}
	for {
		payDetailID := idGeneratorPaymentDetail.NextID()
		var existDetail models.Payment_Details
		scan := server.DB.Where("pay_detail_id = ?", payDetailID).First(&existDetail)
		if scan.RowsAffected == 0 {
			pay_detail = &models.Payment_Details{
				Pay_Detail_ID:  payDetailID,
				Installment_ID: installmentID,
				Pay_Date:       pay_date,
				Pay_Amount:     models.Decimal{Decimal: pay_amount},
				Created_At:     time.Now(),
				Updated_At:     time.Now(),
			}
			break
		}
	}

	pDetails := models.Payment_Details{}
	payDetailShow, err := pDetails.CreatePaymentDetails(server.DB, pay_detail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var paymentStatus models.Payment_Status
	err = server.DB.Preload("Adjustments").First(&paymentStatus, "invoice_id = ?", installment.Invoice_ID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var total decimal.Decimal
	if len(paymentStatus.Adjustment) == 0 {
		total = installment.Ins_Amount.Decimal // Mengambil nilai Ins_Amount dari installment
	} else {
		for _, adj := range paymentStatus.Adjustment {
			total = total.Add(adj.Adjustment_Amount.Decimal) // Mengambil AdjustmentAmount dari setiap Adjustment dan menambahkannya ke total
		}
	}

	paymentBalance := decimal.Zero
	if prevBalance.LessThan(decimal.Zero) {
		paymentBalance = paymentAllocation.Add(prevBalance)
	} else {
		paymentBalance = paymentAllocation.Sub(total)
	}

	response := map[string]interface{}{
		"payment_detail":  payDetailShow,
		"total":           total,
		"payment_balance": paymentBalance,
	}
	json.NewEncoder(w).Encode(response)
}
