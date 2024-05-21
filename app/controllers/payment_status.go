package controllers

import (
	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
)

func (server *Server) CreateNewPaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	payment_S_ID := vars["payment_status_id"]

	invoiceModel := models.Invoice{}
	invoices, err := invoiceModel.GetInvoice(server.DB)
	if err != nil {
		http.Error(w, "Retrive Invoice fail - payment", http.StatusBadRequest)
		return 
	}

	


}