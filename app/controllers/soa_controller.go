package controllers

import (
	"encoding/json"
	
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	
	"github.com/gorilla/mux"
	
)


func (server *Server) CreateSoaAction(w http.ResponseWriter, r *http.Request) {

	name_of_insured_soa := r.FormValue("name_of_insured_soa")
	periode_start_soa := r.FormValue("periode_start_soa")
	periode_end_soa := r.FormValue("periode_end_soa")

	const layoutTime = "2006-01-02"

	if name_of_insured_soa == "" || periode_start_soa == "" || periode_end_soa == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
	}

	p_start_soa, err := time.Parse(layoutTime, periode_start_soa)
	if err != nil {
		http.Error(w, "invalid format!"+err.Error(), http.StatusBadRequest)
		return
	}
	
	p_end_soa, err := time.Parse(layoutTime, periode_end_soa)
	if err != nil {
		http.Error(w, "invalid format!"+err.Error(), http.StatusBadRequest)
		return
	}

	soa_M := models.Statement_Of_Account{}
	var idGeneratorSOA = NewIDGenerator("SOA")

	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	soa := &models.Statement_Of_Account{}
	for {
		soaID := idGeneratorSOA.NextID()

		var existSOA models.Statement_Of_Account

		scan := server.DB.Where("soa_id = ?", soaID).First(&existSOA)

		if scan.RowsAffected == 0 {
			soa = &models.Statement_Of_Account{
				SOA_ID:          soaID,
				UserID:          userID,
				Name_Of_Insured: name_of_insured_soa,
				Period_Start:    p_start_soa,
				Period_End:      p_end_soa,
				Created_At:      time.Now(),
				Updated_At:      time.Now(),
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


func (server *Server) DeleteSoaAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	soa_id := vars["soa_id"]

	soa := &models.Statement_Of_Account{}
	err := soa.DeleteSOA(server.DB, soa_id)
	if err != nil {
		http.Error(w, "Delete fail!"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SOA deleted successfully"))
}

/*

currentDate := time.Now()
    aging := int(payment.PaymentDate.Sub(currentDate).Hours() / 24)

    var message string
    if aging < 0 {
        message = fmt.Sprintf("Sudah berlalu lebih dari %d hari", -aging)
    } else {
        message = fmt.Sprintf("Aging: %d hari", aging)
    }

*/
