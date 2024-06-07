package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
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


func (server *Server) GetItemsBySoaID(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	soa_id_controller := vars["soa_id"]

	var items models.Statement_Of_Account_Details
	itemsShow, err := items.GetItemsBySoaID(server.DB, soa_id_controller)
	if err != nil{
		http.Error(w, "retrive items by soa_id fail!"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemsShow)
}