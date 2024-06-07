package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	
)


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