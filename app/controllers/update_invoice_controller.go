package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	
	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	
	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	
)

func (server *Server) UpdateInvoices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoice_id"]

	const layoutTime = "2006-01-02"

	
	recipient := r.FormValue("recipient")
	address := r.FormValue("address")
	net_premium := r.FormValue("net_premium")
	desc_discount := r.FormValue("desc_discount")
	desc_admin_cost := r.FormValue("desc_admin_cost")
	desc_risk_management := r.FormValue("desc_risk_management")
	desc_brokage := r.FormValue("desc_brokage")
	desc_pph := r.FormValue("desc_pph")
	total_premium_due := r.FormValue("total_premium_due")
	policy_number := r.FormValue("policy_number")
	name_of_insured := r.FormValue("name_of_insured")
	address_of_insured := r.FormValue("address_of_insured")
	type_of_insurance := r.FormValue("type_of_insurance")
	periode_start := r.FormValue("periode_start")
	periode_end := r.FormValue("periode_end")
	terms_of_period := r.FormValue("terms_of_period")
	remarks := r.FormValue("remarks")

	if recipient == "" || address == "" || net_premium == "" ||
		desc_discount == "" || desc_admin_cost == "" || desc_risk_management == "" ||
		desc_brokage == "" || desc_pph == "" || total_premium_due == "" || policy_number == "" ||
		name_of_insured == "" || address_of_insured == "" || type_of_insurance == "" ||
		periode_start == "" || periode_end == "" || terms_of_period == "" || remarks == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}

	p_Start, err := time.Parse(layoutTime, periode_start)
	if err != nil {
		http.Error(w, "invalid periode start!", http.StatusBadRequest)
		return
	}
	p_End, err := time.Parse(layoutTime, periode_end)
	if err != nil {
		http.Error(w, "invalid periode end!", http.StatusBadRequest)
		return
	}

	netPremiumDecimal, err := convertToDecimal(net_premium)
	if err != nil {
		http.Error(w, "invalid net premium!", http.StatusBadRequest)
		return
	}
	descDiscountDecimal, err := convertToDecimal(desc_discount)
	if err != nil {
		http.Error(w, "invalid discount!", http.StatusBadRequest)
		return
	}
	descAdminCostDecimal, err := convertToDecimal(desc_admin_cost)
	if err != nil {
		http.Error(w, "invalid admin cost!", http.StatusBadRequest)
		return
	}
	descRiskManagementDecimal, err := convertToDecimal(desc_risk_management)
	if err != nil {
		http.Error(w, "invalid risk management cost!", http.StatusBadRequest)
		return
	}
	descBrokageDecimal, err := convertToDecimal(desc_brokage)
	if err != nil {
		http.Error(w, "invalid brokage cost!", http.StatusBadRequest)
		return
	}
	descPPHDecimal, err := convertToDecimal(desc_pph)
	if err != nil {
		http.Error(w, "invalid PPH cost!", http.StatusBadRequest)
		return
	}
	totalPremiumDueDecimal, err := convertToDecimal(total_premium_due)
	if err != nil {
		http.Error(w, "invalid total premium due!", http.StatusBadRequest)
		return
	}

	/* Installment form */
	due_date := r.FormValue("due_date")
	ins_amount := r.FormValue("ins_amount")

	var installment models.Installment
	installmentV := false
	if due_date == "" || ins_amount == "" {
		d_date, err := time.Parse(layoutTime, due_date)
		if err != nil {
			http.Error(w, "invalid due date!", http.StatusBadRequest)
			return
		}

		insAmountDecimal, err := convertToDecimal(ins_amount)
		if err != nil {
			http.Error(w, "invalid installment amount!", http.StatusBadRequest)
			return
		}

		// Retrieve existing installment by invoice_id
		var existingInstallment models.Installment
		if err := server.DB.Where("invoice_id = ?", invoiceID).First(&existingInstallment).Error; err != nil {
			http.Error(w, "Installment not found!", http.StatusNotFound)
			return
		}

		installment = models.Installment{
			Installment_ID: existingInstallment.Installment_ID,
			Invoice_ID:     invoiceID,
			Due_Date:       d_date,
			Ins_Amount:     insAmountDecimal,
		}
		http.Error(w, "Fill installment", http.StatusBadRequest)
		return 
	}

	/* Sum insured form */
	items_name := r.FormValue("items_name")
	sum_ins_amount := r.FormValue("sum_ins_amount")
	notes := r.FormValue("notes")

	var sumInsuredDetail models.Sum_Insured_Details
	sumInsuredDetailsV := false
	if items_name == "" || sum_ins_amount == "" {
		sumInsAmountDecimal, err := convertToDecimal(sum_ins_amount)
		if err != nil {
			http.Error(w, "invalid sum insured amount!", http.StatusBadRequest)
			return
		}

		// Retrieve existing sum insured detail by invoice_id
		var existingSumInsured models.Sum_Insured_Details
		if err := server.DB.Where("invoice_id = ?", invoiceID).First(&existingSumInsured).Error; err != nil {
			http.Error(w, "Sum insured detail not found!", http.StatusNotFound)
			return
		}

		sumInsuredDetail = models.Sum_Insured_Details{
			Sum_Insured_ID:     existingSumInsured.Sum_Insured_ID,
			Invoice_ID:         invoiceID,
			Items_Name:         items_name,
			Sum_Insured_Amount: sumInsAmountDecimal,
			Notes:              notes,
		}
		http.Error(w, "Fill sum insured", http.StatusBadRequest)
		return 
	}

	/* Get user token */
	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	invoice := &models.Invoice{
		Invoice_ID:           invoiceID,
		UserID:               userID,
		Recipient:            recipient,
		Address:              address,
		Net_Premium:          netPremiumDecimal,
		Desc_Discount:        descDiscountDecimal,
		Desc_Admin_Cost:      descAdminCostDecimal,
		Desc_Risk_Management: descRiskManagementDecimal,
		Desc_Brokage:         descBrokageDecimal,
		Desc_PPH:             descPPHDecimal,
		Total_Premium_Due:    totalPremiumDueDecimal,
		Policy_Number:        policy_number,
		Name_Of_Insured:      name_of_insured,
		Address_Of_Insured:   address_of_insured,
		Type_Of_Insurance:    type_of_insurance,
		Period_Start:         p_Start,
		Period_End:           p_End,
		Terms_Of_Period:      terms_of_period,
		Remarks:              remarks,
		Updated_At:           time.Now(),
	}
	installments := []models.Installment{}
	if installmentV {
		installments = append(installments, installment)
	}

	sumInsuredDetails := []models.Sum_Insured_Details{}
	if sumInsuredDetailsV {
		sumInsuredDetails = append(sumInsuredDetails, sumInsuredDetail)
	}

	// Update the invoice
	err = invoice.UpdateInvoices(server.DB, invoiceID, installments, sumInsuredDetails)
	if err != nil {
		http.Error(w, "Failed to update invoice!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Invoice has just been updated",
		"invoice": invoice,
	}
	json.NewEncoder(w).Encode(response)
}
