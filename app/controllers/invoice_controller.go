package controllers

import (
	"encoding/json"
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

func (server *Server) CreateInvoicesAction(w http.ResponseWriter, r *http.Request){
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

	if typeInvoice == "" || recipient == "" || address == "" || policy_number == "" ||
		name_of_insured == "" || address_of_insured == "" || type_of_insurance == "" ||
		periode_start == "" || periode_end == "" || terms_of_period =="" || remarks == "" {
			http.Error(w, "Please filltherequired fields!", http.StatusSeeOther)
			return
	}

	p_Start, err := time.Parse(time.Stamp, periode_start)
	if err != nil{
		http.Error(w, "invalid periode start!", http.StatusBadRequest)
		return
	}
	p_End, err := time.Parse(time.Stamp, periode_end)
	if err != nil{
		http.Error(w, "invalid periode end!", http.StatusBadRequest)
		return
	}

}
