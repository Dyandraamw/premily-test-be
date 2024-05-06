package controllers

import (
	"encoding/json"
	"net/http"

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
