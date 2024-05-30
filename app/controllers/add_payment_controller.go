package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/shopspring/decimal"
)

func (server *Server) AddPayment(w http.ResponseWriter, r *http.Request) {
	installmentID := r.FormValue("installment_id")
	p_date_detail := r.FormValue("pay_date")
	p_amount_detail := r.FormValue("pay_amount")

	if installmentID == "" || p_date_detail == "" || p_amount_detail == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}

	var installment models.Installment
	err := server.DB.Preload("Payment_Details").First(&installment, "installment_id = ?", installmentID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (decimal.Decimal, error) {
		return decimal.NewFromString(s)
	}

	const layoutTime = "02-01-2006"
	pay_date, err := time.Parse(layoutTime, p_date_detail)
	if err != nil {
		http.Error(w, "invalid payment date!", http.StatusBadRequest)
		return
	}

	pay_amount, err := convertToDecimal(p_amount_detail)
	if err != nil {
		http.Error(w, "invalid payment amount!", http.StatusBadRequest)
		return
	}

	// var prevBalance decimal.Decimal
	// for _, detail := range installment.Payment_Details {
	// 	prevBalance = prevBalance.Add(detail.Pay_Amount.Decimal)
	// }

	// var paymentAllocation decimal.Decimal
	// if prevBalance.LessThan(decimal.Zero) {
	// 	paymentAllocation = pay_amount
	// } else {
	// 	paymentAllocation = prevBalance.Sub(pay_amount)
	// }

	var idGeneratorPaymentDetail = NewIDGenerator("Payment-Details")

	pay_detail := &models.Payment_Details{}
	for {
		payDetailID := idGeneratorPaymentDetail.NextID()
		var existDetail models.Payment_Details
		scan := server.DB.Where("pay_detail_id = ?", payDetailID).First(&existDetail)
		if scan.RowsAffected == 0 {
			pay_detail = &models.Payment_Details{
				Pay_Detail_ID:  payDetailID,
				Installment_ID: installmentID,
				Pay_Date:       pay_date,
				Pay_Amount:     models.Decimal{Decimal: pay_amount},
				Created_At:     time.Now(),
				Updated_At:     time.Now(),
			}
			break
		}
	}

	pDetails := models.Payment_Details{}
	pDetailShow, err := pDetails.CreatePaymentDetails(server.DB, pay_detail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(pDetailShow)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

