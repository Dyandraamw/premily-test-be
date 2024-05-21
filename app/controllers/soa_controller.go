package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func (server *Server) CreateSoaAction(w http.ResponseWriter, r *http.Request) {

	name_of_insured_soa := r.FormValue("name_of_insured_soa")
	periode_start_soa := r.FormValue("periode_start_soa")
	periode_end_soa := r.FormValue("periode_end_soa")

	const layoutTime = "02-01-2006"

	if name_of_insured_soa == "" || periode_start_soa == "" || periode_end_soa == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
	}

	p_start_soa, err := time.Parse(layoutTime, periode_start_soa)
	if err != nil {
		return
	}

	p_end_soa, err := time.Parse(layoutTime, periode_end_soa)
	if err != nil {
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

func (server *Server) AddItemSoaAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SoA_ID := vars["soa_id"]

	invoiceModel := models.Invoice{}
	invoices, err := invoiceModel.GetInvoice(server.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Get invoice fail", http.StatusBadRequest)
		return
	}

	// var idGeneratorSoaDetails = NewIDGenerator("SOA-Item")

	for _, invoice := range *invoices {
		installmentModel := models.Installment{}
		installments_M, err := installmentModel.GetInstallmentByInvoiceID(server.DB, invoice.Invoice_ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			http.Error(w, "Get installments fail", http.StatusBadRequest)
			return
		}
		for x, installment := range *installments_M {
			instalmentStanding := x + 1

			paymentDate := r.FormValue("payment_date")
			paymentAmount := r.FormValue("payment_amount")

			const layoutTime = "02-01-2006"

			if paymentDate == "" || paymentAmount == "" {
				http.Error(w, "Please fill the required fields!", http.StatusBadRequest)
				return
			}

			p_date_soa_details, err := time.Parse(layoutTime, paymentDate)
			if err != nil {
				http.Error(w, "Invalid payment date format", http.StatusBadRequest)
				return
			}

			convertToDecimal := func(s string) (decimal.Decimal, error) {
				d, err := decimal.NewFromString(s)
				if err != nil {
					return decimal.Decimal{}, err
				}
				return d, nil
			}

			p_amount_soa_details, err := convertToDecimal(paymentAmount)
			if err != nil {
				http.Error(w, "Invalid payment amount format", http.StatusBadRequest)
				return
			}

			var status_SOA_Items string
			paymentAllocation := p_amount_soa_details.Sub(installment.Ins_Amount)
			if paymentAllocation.Equal(decimal.Zero){
				status_SOA_Items= "PAID"
			}else if paymentAllocation.IsPositive() {
				status_SOA_Items = "PAID"
			} else {
				status_SOA_Items = "OUTSTANDING"
			}

			currentDate := time.Now()
			AgingDay := int(p_date_soa_details.Sub(currentDate).Hours() / 24)

			if AgingDay < 0 {
				fmt.Printf("Sudah berlalu lebih dari %d hari", -AgingDay)
			} else {
				fmt.Printf("Aging: %d hari", AgingDay)
			}

			soaDetails := &models.Statement_Of_Account_Details{
				SOA_Details_ID:       uuid.New().String(),
				SOA_ID:               SoA_ID,
				Invoice_ID:           invoice.Invoice_ID,
				Recipient:            invoice.Recipient,
				Installment_Standing: uint(instalmentStanding),
				Due_Date:             installment.Due_Date,
				SOA_Amount:           installment.Ins_Amount,
				Payment_Date:         p_date_soa_details,
				Payment_Amount:       p_amount_soa_details,
				Payment_Allocation:   p_amount_soa_details,
				Status:               status_SOA_Items,
				Aging:                uint(AgingDay),
				Created_At:           time.Now(),
				Updated_At:           time.Now(),
			}

			_, err = soaDetails.CreateSoaDetails(server.DB, soaDetails)
			if err != nil {
				http.Error(w, "Create items fail!", http.StatusBadRequest)
				return
			}

		}
		
	}
	var result models.Invoice
	err = server.DB.Preload("Installment").Preload("Sum_Insured_Details").Where("invoice_id = ?",(*invoices)[0].Invoice_ID).First(&result).Error
	if err != nil {
		http.Error(w, "Failed to retrieve updated invoice", http.StatusInternalServerError)
		return
	}

	// Marshal the result into JSON
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("SOA details added successfully"))

}

func (server *Server) DeleteSoaAction(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	soa_id := vars["soa_id"]

	soa := &models.Statement_Of_Account{}
	err := soa.DeleteSOA(server.DB, soa_id)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
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
