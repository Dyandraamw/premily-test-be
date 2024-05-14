package controllers

import (
	"encoding/json"
	_ "math/rand"
	"net/http"
	"time"

	"fmt"

	"github.com/frangklynndruru/premily_backend/app/models"
)

func (server *Server) Invoice(w http.ResponseWriter, r *http.Request) {
	invoiceModel := models.Invoice{}

	res, err := invoiceModel.GetInvoice(server.DB)
	fmt.Println(res)
	// fmt.Println(invoiceModel.GetInvoice(server.DB))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (server *Server) CreateInvoicesAction(w http.ResponseWriter, r *http.Request) {
	typeInvoice := r.FormValue("typeInvoice")
	recipient := r.FormValue("recipient")
	address := r.FormValue("address")
	policy_number := r.FormValue("policy_number")
	name_of_insured := r.FormValue("name_of_insured")
	address_of_insured := r.FormValue("address_of_insured")
	type_of_insurance := r.FormValue("type_of_insurance")
	periode_start := r.FormValue("periode_start")
	periode_end := r.FormValue("periode_end")
	terms_of_period := r.FormValue("terms_of_period")
	remarks := r.FormValue("remarks")

	var err error
	if typeInvoice == "" || recipient == "" || address == "" || policy_number == "" ||
		name_of_insured == "" || address_of_insured == "" || type_of_insurance == "" ||
		periode_start == "" || periode_end == "" || terms_of_period == "" || remarks == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		
		return
	}
	const layout ="02-01-2006"
	p_Start, err := time.Parse(layout, periode_start)
	if err != nil {
		http.Error(w, "invalid periode start!", http.StatusBadRequest)
		return
	}
	p_End, err := time.Parse(layout, periode_end)
	if err != nil {
		http.Error(w, "invalid periode end!", http.StatusBadRequest)
		return
	}

	invoices_M := models.Invoice{}
	invoice_ID_Generate, err := models.GenerateInvoiceID(server.DB, models.Type(typeInvoice))
	if err != nil {
		http.Error(w, "Create invoices fail - ID Fail!", http.StatusSeeOther)
	}
	invoices := &models.Invoice{
		Invoice_ID:         invoice_ID_Generate,
		Type:               models.Type(typeInvoice),
		Address:            address,
		Policy_Number:      policy_number,
		Name_Of_Insured:    name_of_insured,
		Address_Of_Insured: address_of_insured,
		Type_Of_Insurance:  type_of_insurance,
		Period_Start:       p_Start,
		Period_End:         p_End,
		Terms_Of_Period:    terms_of_period,
		Remarks:            remarks,
	}

	invoicesModels, err := invoices_M.CreateInvoices(server.DB, invoices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Create invoices fail!")
		return
	}

	data, _ := json.Marshal(invoicesModels)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
