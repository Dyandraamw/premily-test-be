package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
)

func (server *Server) GetSoaResponseList(w http.ResponseWriter, r *http.Request) {
	soaModel := models.Statement_Of_Account{}

	soas, err := soaModel.GetSoaList(server.DB)
	if err != nil {
		http.Error(w, "Retrive statement of acoount fail!"+err.Error(), http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(soas)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
