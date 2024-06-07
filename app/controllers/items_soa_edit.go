package controllers

import (
	"encoding/json"

	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	
	"github.com/gorilla/mux"
	
)

func (server *Server) UpdateItemSoaAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SoaDet_ID := vars["soa_details_id"]

	// Fetch existing SOA details
	var existingSOADetails models.Statement_Of_Account_Details
	if err := server.DB.First(&existingSOADetails, "soa_details_id = ?", SoaDet_ID).Error; err != nil {
		http.Error(w, "SOA Details not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Parse input values
	recipient := r.FormValue("recipient")
	due_date := r.FormValue("due_date")
	soa_amount := r.FormValue("soa_amount")
	payment_date := r.FormValue("payment_date")
	payment_amount := r.FormValue("payment_amount")

	// Validate non-empty fields
	if recipient == "" || due_date == "" || soa_amount == "" || payment_date == "" || payment_amount == "" {
		http.Error(w, "Fill the fields!", http.StatusBadRequest)
		return
	}

	const layoutTime = "2006-01-02"

	// Parse and validate due_date
	dueDate, err := time.Parse(layoutTime, due_date)
	if err != nil {
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	// Convert and validate soa_amount
	soaAmount, err := convertToDecimal(soa_amount)
	if err != nil {
		http.Error(w, "Invalid SOA amount format", http.StatusBadRequest)
		return
	}

	// Parse and validate payment_date
	paymentDate, err := time.Parse(layoutTime, payment_date)
	if err != nil {
		http.Error(w, "Invalid payment date format", http.StatusBadRequest)
		return
	}

	// Convert and validate payment_amount
	paymentAmount, err := convertToDecimal(payment_amount)
	if err != nil {
		http.Error(w, "Invalid payment amount format", http.StatusBadRequest)
		return
	}

	// Update the existing SOA details
	existingSOADetails.Recipient = recipient
	existingSOADetails.Due_Date = dueDate
	existingSOADetails.SOA_Amount = soaAmount
	existingSOADetails.Payment_Date = paymentDate
	existingSOADetails.Payment_Amount = paymentAmount
	existingSOADetails.Updated_At = time.Now()

	// Save the updated SOA details to the database
	if err := existingSOADetails.UpdatesItemsSoa(server.DB, SoaDet_ID); err != nil {
		http.Error(w, "Failed to update SOA Details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated SOA details to include in the response
	if err := server.DB.First(&existingSOADetails, "soa_details_id = ?", SoaDet_ID).Error; err != nil {
		http.Error(w, "Failed to fetch updated SOA Details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the updated SOA details into JSON
	data, err := json.Marshal(existingSOADetails)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}