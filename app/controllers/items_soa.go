package controllers

import (
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)
func (server *Server) AddItemSoaAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SoA_ID := vars["soa_id"]

	if SoA_ID == "" {
		http.Error(w, "SOA_ID is required", http.StatusBadRequest)
		return
	}

	invoiceModel := models.Invoice{}
	invoices, err := invoiceModel.GetInvoice(server.DB)
	if err != nil {
		http.Error(w, "Get invoice fail: "+err.Error(), http.StatusBadRequest)
		return
	}

	var responseSoaDetails []*models.Statement_Of_Account_Details

	for _, invoice := range *invoices {
		installmentModel := models.Installment{}
		installments_M, err := installmentModel.GetInstallmentByInvoiceID(server.DB, invoice.Invoice_ID)
		if err != nil {
			http.Error(w, "Get installments fail: "+err.Error(), http.StatusBadRequest)
			return
		}

		selectedInvoiceID := r.FormValue("invoice_id")
		if selectedInvoiceID == "" {
			http.Error(w, "Please select invoice!", http.StatusBadRequest)
			return
		}

		var selectedInvoice *models.Invoice
		for _, invoiceList := range *invoices {
			if invoiceList.Invoice_ID == selectedInvoiceID {
				selectedInvoice = &invoiceList
				break
			}
		}

		if selectedInvoice == nil {
			http.Error(w, "Selected invoice not found!", http.StatusBadRequest)
			return
		}

		for x, installment := range *installments_M {
			instalmentStanding := x + 1

			paymentDate := r.FormValue("payment_date")
			paymentAmount := r.FormValue("payment_amount")

			const layoutTime = "2006-01-02"

			if paymentDate == "" || paymentAmount == "" {
				http.Error(w, "Please fill the required fields!", http.StatusBadRequest)
				return
			}

			p_date_soa_details, err := time.Parse(layoutTime, paymentDate)
			if err != nil {
				http.Error(w, "Invalid payment date format", http.StatusBadRequest)
				return
			}

			convertToDecimal := func(s string) (models.Decimal, error) {
				d, err := decimal.NewFromString(s)
				if err != nil {
					return models.Decimal{}, err
				}
				return models.Decimal{Decimal: d}, nil
			}

			p_amount_soa_details, err := convertToDecimal(paymentAmount)
			if err != nil {
				http.Error(w, "Invalid payment amount format", http.StatusBadRequest)
				return
			}

			insAmountDecimal, err := decimal.NewFromString(installment.Ins_Amount.String())
			if err != nil {
				http.Error(w, "Invalid insured amount format", http.StatusBadRequest)
				return
			}

			var status_SOA_Items string
			paymentAllocation := p_amount_soa_details.Sub(insAmountDecimal)
			if paymentAllocation.Equal(decimal.Zero) {
				status_SOA_Items = "PAID"
			} else if paymentAllocation.IsPositive() {
				status_SOA_Items = "PAID"
			} else {
				status_SOA_Items = "OUTSTANDING"
			}

			currentDate := time.Now()
			AgingDay := int(p_date_soa_details.Sub(currentDate).Hours() / 24)

			soaDetails := &models.Statement_Of_Account_Details{
				SOA_Details_ID:       uuid.New().String(),
				SOA_ID:               SoA_ID,
				Invoice_ID:           selectedInvoiceID,
				Recipient:            selectedInvoice.Recipient,
				Installment_Standing: uint(instalmentStanding),
				Due_Date:             installment.Due_Date,
				SOA_Amount:           installment.Ins_Amount,
				Payment_Date:         p_date_soa_details,
				Payment_Amount:       p_amount_soa_details,
				Payment_Allocation:   p_amount_soa_details,
				Status:               status_SOA_Items,
				Aging:                uint(AgingDay),
				Created_At:           currentDate,
				Updated_At:           currentDate,
			}

			log.Printf("Creating SOA Details: %+v\n", soaDetails)

			if _, err := soaDetails.CreateSoaDetails(server.DB, soaDetails); err != nil {
				http.Error(w, "Create items fail: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Tambahkan soaDetails ke response list
			responseSoaDetails = append(responseSoaDetails, soaDetails)
		}
	}

	// Marshal the responseSoaDetails into JSON
	data, err := json.Marshal(responseSoaDetails)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}





// Helper function to convert string to Decimal
func convertToDecimal(s string) (models.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return models.Decimal{}, err
	}
	return models.Decimal{Decimal: d}, nil
}
