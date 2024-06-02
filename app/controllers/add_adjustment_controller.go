package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func (server *Server) AddAjustment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pStatusID := vars["payment_status_id"]
	adjustment_title := r.FormValue("title")
	adjustment_amount := r.Form["amount"]

	var err error
	if adjustment_title == "" || len(adjustment_amount) == 0 {
		// http.Error(w, err.Error(), http.StatusSeeOther)
		fmt.Println("ini panjang adjsutment : ",len(adjustment_amount))
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (decimal.Decimal, error) {
		return decimal.NewFromString(s)
	}

	var adjust_amounts []decimal.Decimal

	for _, amt := range adjustment_amount {

		adjust_amount, err := convertToDecimal(amt)
		if err != nil {
			http.Error(w, "invalid payment amount!", http.StatusBadRequest)
			return
		}
		adjust_amounts = append(adjust_amounts, adjust_amount)
	}

	// Get the invoice_id from the payment_status_id
	var paymentStatus models.Payment_Status
	err = server.DB.Where("payment_status_id = ?", pStatusID).First(&paymentStatus).Error
	if err != nil {
		http.Error(w, "Error fetching payment status!", http.StatusInternalServerError)
		return
	}

	invoiceID := paymentStatus.Invoice_ID

	// Get installments related to the invoice_id
	var installments []models.Installment
	err = server.DB.Where("invoice_id = ?", invoiceID).Find(&installments).Error
	if err != nil {
		http.Error(w, "Error fetching installments!", http.StatusInternalServerError)
		return
	}

	// Ensure the number of adjustment amounts matches the number of installments
	if len(installments) != len(adjust_amounts) {
		fmt.Println(len(installments))
		http.Error(w, "Adjustment amounts do not match the number of installments!", http.StatusBadRequest)
		return
	}

	var idGeneratorAdjustment = NewIDGenerator("Adjustment")
	newAdjustments := []*models.Adjustment{}
	adjustM := &models.Adjustment{}
	for x, _ := range installments {
		adjust_id := idGeneratorAdjustment.NextID()
		// var existAdjustment models.Adjustment
		// scan := server.DB.Where("adjustment_id = ?", adjust_id).First(&existAdjustment)
		// if scan.RowsAffected == 0 {
		adjustM = &models.Adjustment{
			Adjustment_ID:     adjust_id,
			Payment_Status_ID: pStatusID,
			Adjustment_Title:  adjustment_title,
			Adjustment_Amount: models.Decimal{Decimal: adjust_amounts[x]},
			Created_At:        time.Now(),
			Updated_At:        time.Now(),
		}

		newAdjustment, err := adjustM.CreateAdjustment(server.DB, adjustM)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newAdjustments = append(newAdjustments, newAdjustment)
		// 	break
		// }
	}
	// Prepare response data

	var installmentPointers []*models.Installment
	for i := range installments {
		installmentPointers = append(installmentPointers, &installments[i])
	}

	responseData := struct {
		Installments []*models.Installment `json:"installments"`
		Adjustments  []*models.Adjustment  `json:"adjustments"`
	}{
		Installments: installmentPointers,
		Adjustments:  newAdjustments,
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

	// Prepare response data
	responseData := struct {
		PaymentStatus  models.Payment_Status  `json:"payment_status"`
		Adjustments    []models.Adjustment    `json:"adjustments"`
		Invoice        models.Invoice         `json:"invoice"`
		Installments   []models.Installment   `json:"installments"`
		PaymentDetails []models.Payment_Details `json:"payment_details"`
	}{
		PaymentStatus:  paymentStatus,
		Adjustments:    paymentStatus.Adjustment,
		Invoice:        invoice,
		Installments:   installments,
		PaymentDetails: paymentDetails,
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
