package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

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

	const layoutTime = "2006-01-02"
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

	// var prevBalance decimal.Decimal
	// for _, detail := range installment.Payment_Details {
	// 	prevBalance = prevBalance.Add(detail.Pay_Amount.Decimal)
	// }

	// var paymentAllocation decimal.Decimal
	// if prevBalance.LessThan(decimal.Zero) {
	// 	paymentAllocation = pay_amount
	// } else {
	// 	paymentAllocation = prevBalance.Sub(pay_amount)
	// }

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
	pDetailShow, err := pDetails.CreatePaymentDetails(server.DB, pay_detail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(pDetailShow)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (server *Server) UpdatePaymentDetails(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	pay_details_id := vars["pay_detail_id"]

	var existPayment models.Payment_Details
	err := server.DB.First(&existPayment, "pay_detail_id = ?", pay_details_id).Error
	if err != nil{
		http.Error(w, "Payment not found!"+err.Error(), http.StatusNotFound)
	}
	const layoutTime = "2006-01-02"

	pay_date := r.FormValue("pay_date")
	pay_amount := r.FormValue("pay_amount")

	if  pay_date == "" || pay_amount == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}
	// Parse and validate pay_date
	pay_date_details, err := time.Parse(layoutTime, pay_date)
	if err != nil {
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}
	pay_amount_details, err := convertToDecimal(pay_amount)
	if err != nil {
		http.Error(w, "Invalid payment amount format", http.StatusBadRequest)
		return
	}

	existPayment.Pay_Date = pay_date_details
	existPayment.Pay_Amount = pay_amount_details
	existPayment.Updated_At = time.Now()

	err = existPayment.UpdatePayment(server.DB, pay_details_id )
	if err != nil {
		http.Error(w, "Failed to update payment details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = server.DB.First(&existPayment, "pay_detail_id = ?", pay_details_id).Error
	if err != nil{
		http.Error(w, "Response fail!: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the updated SOA details into JSON
	data, err := json.Marshal(existPayment)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}