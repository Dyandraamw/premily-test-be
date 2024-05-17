package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
)

func (server *Server) CreateSoaAction(w http.ResponseWriter, r *http.Request) {

	name_of_insured_soa := r.FormValue("name_of_insured_soa")
	periode_start_soa := r.FormValue("periode_start_soa")
	periode_end_soa := r.FormValue("periode_end_soa")

	const layoutTime = "02-01-2006"

	if name_of_insured_soa == "" || periode_start_soa == "" || periode_end_soa ==""{
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
	}

	p_start_soa, err := time.Parse(layoutTime, periode_start_soa)
	if err != nil{
		return
	}

	p_end_soa, err := time.Parse(layoutTime, periode_end_soa)
	if err != nil{
		return
	}

	soa_M := models.Statement_Of_Account{}
	var idGeneratorSOA = NewIDGenerator("SOA")

	userID, err:= auth.GetTokenUserLogin(r.Context())
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	soa := &models.Statement_Of_Account{}
	for{
		soaID := idGeneratorSOA.NextID()

		var existSOA models.Statement_Of_Account

		scan := server.DB.Where("soa_id = ?", soaID).First(&existSOA)

		if scan.RowsAffected == 0{
			soa = &models.Statement_Of_Account{
				SOA_ID: soaID,
				UserID: userID,
				Name_Of_Insured: name_of_insured_soa,
				Period_Start: p_start_soa,
				Period_End: p_end_soa,
			}
			break
		}
	}

	soaShow, err := soa_M.CreateNewSOA(server.DB, soa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(soaShow)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}


