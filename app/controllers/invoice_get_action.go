package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	
)


func (server *Server) Invoice(w http.ResponseWriter, r *http.Request) {
	invoiceModel := models.Invoice{}

	invoices, err := invoiceModel.GetInvoiceResponseList(server.DB)

	// fmt.Println(invoiceModel.GetInvoice(server.DB))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(invoices)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (server *Server) GetInvoiceByID( w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoice_ID := vars["invoices_id"]

	invoiceModel := models.Invoice{}

	invoices, err := invoiceModel.GetInvoiceByIDmodel(server.DB, invoice_ID)
	if err != nil {
		
		http.Error(w, "Invoice not found!"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}
