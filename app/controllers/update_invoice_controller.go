package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	
	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func (server Server) UpdateInvoices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoice_id"]

	var err error
	const layoutTime = "2006-01-02"

	typeInvoice := r.FormValue("typeInvoice")
	recipient := r.FormValue("recipient")
	address := r.FormValue("address")

	desc_premium := r.FormValue("desc_premium")
	desc_discount := r.FormValue("desc_discount")
	desc_admin_cost := r.FormValue("desc_admin_cost")
	desc_risk_management := r.FormValue("desc_risk_management")
	desc_brokage := r.FormValue("desc_brokage")
	desc_pph := r.FormValue("desc_pph")

	policy_number := r.FormValue("policy_number")
	name_of_insured := r.FormValue("name_of_insured")
	address_of_insured := r.FormValue("address_of_insured")
	type_of_insurance := r.FormValue("type_of_insurance")
	periode_start := r.FormValue("periode_start")
	periode_end := r.FormValue("periode_end")
	terms_of_period := r.FormValue("terms_of_period")
	remarks := r.FormValue("remarks")

	//installment form
	// due_date := r.FormValue("due_date")
	// ins_amount := r.FormValue("ins_amount")

	// cum insured details form
	// items_name := r.FormValue("items_name")
	// sum_ins_amount := r.FormValue("sum_ins_amount")
	// notes := r.FormValue("notes")

	// d_date, err := time.Parse(layoutTime, due_date)
	// if err != nil {
	// 	http.Error(w, "invalid periode start!", http.StatusBadRequest)
	// 	return
	// }

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

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (models.Decimal, error) {
		d, err := decimal.NewFromString(s)
		if err != nil {
			return models.Decimal{}, err
		}
		return models.Decimal{Decimal: d}, nil
	}

	descPremiumDecimal, err := convertToDecimal(desc_premium)
	if err != nil {
		http.Error(w, "invalid description premium!", http.StatusBadRequest)
		return
	}
	descDiscountDecimal, err := convertToDecimal(desc_discount)
	if err != nil {
		http.Error(w, "invalid description discount!", http.StatusBadRequest)
		return
	}
	descAdminCostDecimal, err := convertToDecimal(desc_admin_cost)
	if err != nil {
		http.Error(w, "invalid description admin cost!", http.StatusBadRequest)
		return
	}
	descRiskManagementDecimal, err := convertToDecimal(desc_risk_management)
	if err != nil {
		http.Error(w, "invalid description risk management!", http.StatusBadRequest)
		return
	}
	descBrokageDecimal, err := convertToDecimal(desc_brokage)
	if err != nil {
		http.Error(w, "invalid description brokage!", http.StatusBadRequest)
		return
	}
	descPPHDecimal, err := convertToDecimal(desc_pph)
	if err != nil {
		http.Error(w, "invalid description pph!", http.StatusBadRequest)
		return
	}

	// insAmountDecimal, err := convertToDecimal(ins_amount)
	// if err != nil {
	// 	http.Error(w, "invalid insured amount!", http.StatusBadRequest)
	// 	return
	// }

	// sumInsAmountDecimal, err := convertToDecimal(sum_ins_amount)
	// if err != nil {
	// 	http.Error(w, "invalid insured amount!", http.StatusBadRequest)
	// 	return
	// }

	/*
			GET TOKEN USER LOGIN
		=================================
	*/
	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	invoices_M := models.Invoice{}
	// invoice_ID_Generate, err := models.GenerateInvoiceID(server.DB, models.Type(typeInvoice))
	// if err != nil {
	// 	http.Error(w, "Create invoices fail - ID Fail!", http.StatusSeeOther)
	// }
	invoices := &models.Invoice{
		Invoice_ID:           invoiceID,
		UserID:               userID,
		Type:                 models.Type(typeInvoice),
		Recipient:            recipient,
		Address:              address,
		Net_Premium:          descPremiumDecimal,
		Desc_Discount:        descDiscountDecimal,
		Desc_Admin_Cost:      descAdminCostDecimal,
		Desc_Risk_Management: descRiskManagementDecimal,
		Desc_Brokage:         descBrokageDecimal,
		Desc_PPH:             descPPHDecimal,
		Policy_Number:        policy_number,
		Name_Of_Insured:      name_of_insured,
		Address_Of_Insured:   address_of_insured,
		Type_Of_Insurance:    type_of_insurance,
		Period_Start:         p_Start,
		Period_End:           p_End,
		Terms_Of_Period:      terms_of_period,
		Remarks:              remarks,
	}
	installments := []models.Installment{}
	sum_insured := []models.Sum_Insured_Details{}

	err = json.Unmarshal([]byte(r.FormValue("installments")), &installments)
	if err != nil {
		http.Error(w, "Parsing installments fail!", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(r.FormValue("sum_insured")), &sum_insured)
	if err != nil {
		http.Error(w, "Parsing sum insured fail1", http.StatusBadRequest)
		return
	}

	err = invoices_M.UpdateInvoices(server.DB, invoiceID, installments, sum_insured)
	if err != nil {

		http.Error(w, "update invoice fail"+err.Error(), http.StatusBadRequest)
		return
	}
	

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoices)

}