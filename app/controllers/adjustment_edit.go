package controllers

import(
	"net/http"
	"time"
	"encoding/json"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"

)
func (server *Server) EditAdjustment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adjust_id := vars["adjustment_id"]

	// Validasi input
	adjustment_title := r.FormValue("title")
	adjustment_amount := r.Form["amount"]

	if adjustment_title == "" || len(adjustment_amount) == 0 {
		http.Error(w, "Please fill the required fields!", http.StatusBadRequest)
		return
	}

	// Konversi string ke decimal.Decimal
	convertToDecimal := func(s string) (decimal.Decimal, error) {
		return decimal.NewFromString(s)
	}

	var adjust_amounts []decimal.Decimal

	for _, amt := range adjustment_amount {
		adjust_amount, err := convertToDecimal(amt)
		if err != nil {
			http.Error(w, "Invalid adjustment amount!", http.StatusBadRequest)
			return
		}
		adjust_amounts = append(adjust_amounts, adjust_amount)
	}

	// Dapatkan adjustment yang ada
	var existAdjustment models.Adjustment
	if err := server.DB.First(&existAdjustment, "adjustment_id = ?", adjust_id).Error; err != nil {
		http.Error(w, "Adjustment not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Update nilai adjustment
	existAdjustment.Adjustment_Title = adjustment_title
	existAdjustment.Adjustment_Amount = models.Decimal{Decimal: adjust_amounts[0]}
	existAdjustment.Updated_At = time.Now()

	// Simpan perubahan ke database
	if err := existAdjustment.UpdateAdjustment(server.DB, adjust_id); err != nil {
		http.Error(w, "Failed to update adjustment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil data terbaru dari database untuk dikembalikan sebagai respons
	if err := server.DB.First(&existAdjustment, "adjustment_id = ?", adjust_id).Error; err != nil {
		http.Error(w, "Failed to fetch updated adjustment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal data response ke JSON
	data, err := json.Marshal(existAdjustment)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header dan kirim respons
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}