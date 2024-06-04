package controllers

import (
	"net/http"

	"time"

	"encoding/json"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	
)
func (server *Server) GetPaymentData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pStatusID := vars["payment_status_id"]

	var paymentStatus models.Payment_Status
	err := server.DB.Preload("Adjustment").Where("payment_status_id = ?", pStatusID).First(&paymentStatus).Error
	if err != nil {
		http.Error(w, "Error fetching payment status!", http.StatusInternalServerError)
		return
	}

	var invoice models.Invoice
	err = server.DB.Where("invoice_id = ?", paymentStatus.Invoice_ID).First(&invoice).Error
	if err != nil {
		http.Error(w, "Error fetching invoice!", http.StatusInternalServerError)
		return
	}

	var installments []models.Installment
	err = server.DB.Where("invoice_id = ?", invoice.Invoice_ID).Find(&installments).Error
	if err != nil {
		http.Error(w, "Error fetching installments!", http.StatusInternalServerError)
		return
	}

	var paymentDetails []models.Payment_Details
	for _, installment := range installments {
		var pd []models.Payment_Details
		err = server.DB.Where("installment_id = ?", installment.Installment_ID).Find(&pd).Error
		if err != nil {
			http.Error(w, "Error fetching payment details!", http.StatusInternalServerError)
			return
		}
		paymentDetails = append(paymentDetails, pd...)
	}

	totalAdj, err := models.CalculateAdjustment(server.DB, pStatusID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var installment models.Installment
	if len(installments) > 0 {
		installment = installments[0]  // Assuming you need the first installment
	} else {
		http.Error(w, "No installments found!", http.StatusBadRequest)
		return
	}

	// var installment models.Installment
	total := installment.Ins_Amount.Sub(totalAdj.Decimal)
	if total.IsNegative(){
		http.Error(w, "Calculate premium inception fail!", http.StatusBadRequest)
		return
	}

	totalSum := total.IntPart()

	balance, err := models.CalculatePayment(server.DB, pStatusID, int(totalSum))
	if err != nil {
		http.Error(w, "Calculation balance fail!", http.StatusBadRequest)
		return
	}


	// Prepare response data
	responseData := ResponseData{
		PaymentStatus: paymentStatus.Payment_Status_ID,
		Adjustments:   make([]models.Decimal, len(paymentStatus.Adjustment)),
		Invoice: InvoiceData{
			NameOfInsured: invoice.Name_Of_Insured,
			PeriodStart:   invoice.Period_Start,
			PeriodEnd:     invoice.Period_End,
		},
		Installments:   make([]InstallmentData, len(installments)),
		Payment_Details: make([]PaymentDetailsData, len(paymentDetails)),
		Total: models.Decimal{Decimal: total},
		Balance: models.Decimal{Decimal: balance.Decimal},
		
	}

	// Fill data Adjustments
	for i, adj := range paymentStatus.Adjustment {
		responseData.Adjustments[i] = adj.Adjustment_Amount
	}

	// Fill data Installments
	for i, ins := range installments {
		responseData.Installments[i] = InstallmentData{
			InstallmentID: ins.Installment_ID,
			DueDate:       ins.Due_Date,
			InsAmount:     ins.Ins_Amount,
		}
	}

	// Fill data PaymentDetails
	for i, pd := range paymentDetails {
		responseData.Payment_Details[i] = PaymentDetailsData{
			PayDate:   pd.Pay_Date,
			PayAmount: pd.Pay_Amount,
		}
	}

	// Marshal data response to JSON
	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}

	// Set header and send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}


func (server *Server) PaymentCalculation(w http.ResponseWriter, r *http.Request)  {
	
}



// Struct untuk response data yang lebih spesifik
type ResponseData struct {
	PaymentStatus   string               `json:"payment_status_id"`
	Adjustments     []models.Decimal     `json:"adjustment_amount"`
	Invoice         InvoiceData          `json:"invoice"`
	Installments    []InstallmentData    `json:"installments"`
	Payment_Details []PaymentDetailsData `json:"payment_details"`
	Total	models.Decimal	`json:"total"`
	Balance models.Decimal	`json:"balance"`
}

type InvoiceData struct {
	NameOfInsured string    `json:"name_of_insured"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
}

type InstallmentData struct {
	InstallmentID string         `json:"installment_id"`
	DueDate       time.Time      `json:"due_date"`
	InsAmount     models.Decimal `json:"ins_amount"`
}

type PaymentDetailsData struct {
	PayDate   time.Time      `json:"pay_date"`
	PayAmount models.Decimal `json:"pay_amount"`
}