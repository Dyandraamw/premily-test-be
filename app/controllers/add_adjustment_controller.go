package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func (server *Server) AddAjustment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars()
	pStatusID := vars["payment_status_id"]
	adjustment_title := r.FormValue("title")
	adjustment_amount := r.FormValue("amount")

	if adjustment_title == "" || adjustment_amount == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (decimal.Decimal, error) {
		return decimal.NewFromString(s)
	}

	adjust_amount, err := convertToDecimal(adjustment_amount)
	if err != nil {
		http.Error(w, "invalid payment amount!", http.StatusBadRequest)
		return
	}

	var idGeneratorAdjustment = NewIDGenerator("Adjustment")
	adjustM := &models.Adjustment{}
	for {
		adjust_id := idGeneratorAdjustment.NextID()
		var existAdjustment models.Adjustment
		scan := server.DB.Where("adjustment_id = ?", adjust_id).First(&existAdjustment)
		if scan.RowsAffected == 0 {
			adjustM = &models.Adjustment{
				Adjustment_ID:     adjust_id,
				Payment_Status_ID: pStatusID,
				Adjustment_Title:  adjustment_title,
				Adjustment_Amount: models.Decimal{Decimal: adjust_amount},
				Created_At:        time.Now(),
				Updated_At:        time.Now(),
			}
			break
		}
	}
	adjustmentPoint := models.Adjustment{}
	adjustShow, err := adjustmentPoint.CreateAdjustment(server.DB,adjustM)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(adjustShow)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
