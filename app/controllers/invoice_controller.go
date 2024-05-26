package controllers

import (
	"encoding/json"
	"fmt"
	_ "math/rand"
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/shopspring/decimal"
)

func (server *Server) Invoice(w http.ResponseWriter, r *http.Request) {
	invoiceModel := models.Invoice{}

	invoices, err := invoiceModel.GetInvoice(server.DB)

	// fmt.Println(invoiceModel.GetInvoice(server.DB))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(invoices)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (server *Server) GetInvoiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoice_ID := vars["invoices_id"]

	invoiceModel := models.Invoice{}

	invoices, err := invoiceModel.GetInvoiceByIDmodel(server.DB, invoice_ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Invoice not found!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (server *Server) CreateInvoicesAction(w http.ResponseWriter, r *http.Request) {

	var err error
	const layoutTime = "02-01-2006"

	typeInvoice := r.FormValue("typeInvoice")
	recipient := r.FormValue("recipient")
	address := r.FormValue("address")

	desc_premium := r.FormValue("desc_premium")
	desc_discount := r.FormValue("desc_discount")
	desc_admin_cost := r.FormValue("desc_admin_cost")
	desc_risk_management := r.FormValue("desc_risk_management")
	desc_brokage := r.FormValue("desc_brokage")
	desc_pph := r.FormValue("desc_pph")

	policy_number := r.FormValue("policy_number")
	name_of_insured := r.FormValue("name_of_insured")
	address_of_insured := r.FormValue("address_of_insured")
	type_of_insurance := r.FormValue("type_of_insurance")
	periode_start := r.FormValue("periode_start")
	periode_end := r.FormValue("periode_end")
	terms_of_period := r.FormValue("terms_of_period")
	remarks := r.FormValue("remarks")

	//installment form
	due_date := r.FormValue("due_date")
	ins_amount := r.FormValue("ins_amount")

	//cum insured details form
	items_name := r.FormValue("items_name")
	sum_ins_amount := r.FormValue("sum_ins_amount")
	notes := r.FormValue("notes")

	d_date, err := time.Parse(layoutTime, due_date)
	if err != nil {
		http.Error(w, "invalid periode start!", http.StatusBadRequest)
		return
	}

	if typeInvoice == "" || recipient == "" || address == "" || desc_premium == "" || desc_discount == "" ||
		desc_admin_cost == "" || desc_risk_management == "" || desc_brokage == "" || desc_pph == "" || policy_number == "" ||
		name_of_insured == "" || address_of_insured == "" || type_of_insurance == "" ||
		periode_start == "" || periode_end == "" || terms_of_period == "" || remarks == "" || due_date == "" || ins_amount == "" ||
		items_name == "" || sum_ins_amount == "" || notes == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)

		return
	}

	p_Start, err := time.Parse(layoutTime, periode_start)
	if err != nil {
		http.Error(w, "invalid periode start!", http.StatusBadRequest)
		return
	}
	p_End, err := time.Parse(layoutTime, periode_end)
	if err != nil {
		http.Error(w, "invalid periode end!", http.StatusBadRequest)
		return
	}

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (models.Decimal, error) {
		d, err := decimal.NewFromString(s)
		if err != nil {
			return models.Decimal{}, err
		}
		return models.Decimal{Decimal: d}, nil
	}

	descPremiumDecimal, err := convertToDecimal(desc_premium)
	if err != nil {
		http.Error(w, "invalid description premium!", http.StatusBadRequest)
		return
	}
	descDiscountDecimal, err := convertToDecimal(desc_discount)
	if err != nil {
		http.Error(w, "invalid description discount!", http.StatusBadRequest)
		return
	}
	descAdminCostDecimal, err := convertToDecimal(desc_admin_cost)
	if err != nil {
		http.Error(w, "invalid description admin cost!", http.StatusBadRequest)
		return
	}
	descRiskManagementDecimal, err := convertToDecimal(desc_risk_management)
	if err != nil {
		http.Error(w, "invalid description risk management!", http.StatusBadRequest)
		return
	}
	descBrokageDecimal, err := convertToDecimal(desc_brokage)
	if err != nil {
		http.Error(w, "invalid description brokage!", http.StatusBadRequest)
		return
	}
	descPPHDecimal, err := convertToDecimal(desc_pph)
	if err != nil {
		http.Error(w, "invalid description pph!", http.StatusBadRequest)
		return
	}

	insAmountDecimal, err := convertToDecimal(ins_amount)
	if err != nil {
		http.Error(w, "invalid insured amount!", http.StatusBadRequest)
		return
	}

	sumInsAmountDecimal, err := convertToDecimal(sum_ins_amount)
	if err != nil {
		http.Error(w, "invalid insured amount!", http.StatusBadRequest)
		return
	}

	/*
			GET TOKEN USER LOGIN
		=================================
	*/
	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	invoices_M := models.Invoice{}
	invoice_ID_Generate, err := models.GenerateInvoiceID(server.DB, models.Type(typeInvoice))
	if err != nil {
		http.Error(w, "Create invoices fail - ID Fail!", http.StatusSeeOther)
	}
	invoices := &models.Invoice{
		Invoice_ID:           invoice_ID_Generate,
		UserID:               userID,
		Type:                 models.Type(typeInvoice),
		Recipient:            recipient,
		Address:              address,
		Desc_Premium:         descPremiumDecimal,
		Desc_Discount:        descDiscountDecimal,
		Desc_Admin_Cost:      descAdminCostDecimal,
		Desc_Risk_Management: descRiskManagementDecimal,
		Desc_Brokage:         descBrokageDecimal,
		Desc_PPH:             descPPHDecimal,
		Policy_Number:        policy_number,
		Name_Of_Insured:      name_of_insured,
		Address_Of_Insured:   address_of_insured,
		Type_Of_Insurance:    type_of_insurance,
		Period_Start:         p_Start,
		Period_End:           p_End,
		Terms_Of_Period:      terms_of_period,
		Remarks:              remarks,
		Created_At:           time.Now(),
		Updated_At:           time.Now(),
	}

	_, err = invoices_M.CreateInvoices(server.DB, invoices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Create invoices fail!")
		return
	}

	installment_M := models.Installment{}
	var idGeneratorInstallment = NewIDGenerator("INS")

	installments := &models.Installment{}
	for {
		insID := idGeneratorInstallment.NextID()
		var existInstallment models.Installment

		scan := server.DB.Where("installment_id = ?", insID).First(&existInstallment)
		if scan.RowsAffected == 0 {
			installments = &models.Installment{
				Installment_ID: insID,
				Invoice_ID:     invoices.Invoice_ID,
				Due_Date:       d_date,
				Ins_Amount:     insAmountDecimal,
			}
			break
		}
	}
	_, err = installment_M.CreateInstallment(server.DB, installments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Installment fail!")
		return
	}

	sumIns_M := models.Sum_Insured_Details{}
	var idGeneratorSumIns = NewIDGenerator("S-INS")

	sum_insureds := &models.Sum_Insured_Details{}
	for {
		sum_ins_ID := idGeneratorSumIns.NextID()

		var existSumIns models.Sum_Insured_Details

		scan := server.DB.Where("sum_insured_id = ?", sum_ins_ID).First(&existSumIns)

		if scan.RowsAffected == 0 {
			sum_insureds = &models.Sum_Insured_Details{
				Sum_Insured_ID:     sum_ins_ID,
				Invoice_ID:         invoices.Invoice_ID,
				Items_Name:         items_name,
				Sum_Insured_Amount: sumInsAmountDecimal,
				Notes:              notes,
			}
			break
		}
	}
	_, err = sumIns_M.CreateSumInsuredDetails(server.DB, sum_insureds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Sum insured fail!")
		return
	}

	var result models.Invoice
	err = server.DB.Preload("Installment").Preload("Sum_Insured_Details").Where("invoice_id = ?",
		invoices.Invoice_ID).First(&result).Error

	if err != nil {
		return
	}
	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (server *Server) UpdateInvoices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoice_id"]

	var err error
	const layoutTime = "02-01-2006"

	typeInvoice := r.FormValue("typeInvoice")
	recipient := r.FormValue("recipient")
	address := r.FormValue("address")

	desc_premium := r.FormValue("desc_premium")
	desc_discount := r.FormValue("desc_discount")
	desc_admin_cost := r.FormValue("desc_admin_cost")
	desc_risk_management := r.FormValue("desc_risk_management")
	desc_brokage := r.FormValue("desc_brokage")
	desc_pph := r.FormValue("desc_pph")

	policy_number := r.FormValue("policy_number")
	name_of_insured := r.FormValue("name_of_insured")
	address_of_insured := r.FormValue("address_of_insured")
	type_of_insurance := r.FormValue("type_of_insurance")
	periode_start := r.FormValue("periode_start")
	periode_end := r.FormValue("periode_end")
	terms_of_period := r.FormValue("terms_of_period")
	remarks := r.FormValue("remarks")

	//installment form
	// due_date := r.FormValue("due_date")
	// ins_amount := r.FormValue("ins_amount")

	// cum insured details form
	// items_name := r.FormValue("items_name")
	// sum_ins_amount := r.FormValue("sum_ins_amount")
	// notes := r.FormValue("notes")

	// d_date, err := time.Parse(layoutTime, due_date)
	// if err != nil {
	// 	http.Error(w, "invalid periode start!", http.StatusBadRequest)
	// 	return
	// }

	p_Start, err := time.Parse(layoutTime, periode_start)
	if err != nil {
		http.Error(w, "invalid periode start!", http.StatusBadRequest)
		return
	}
	p_End, err := time.Parse(layoutTime, periode_end)
	if err != nil {
		http.Error(w, "invalid periode end!", http.StatusBadRequest)
		return
	}

	// Convert string to decimal.Decimal
	convertToDecimal := func(s string) (models.Decimal, error) {
		d, err := decimal.NewFromString(s)
		if err != nil {
			return models.Decimal{}, err
		}
		return models.Decimal{Decimal: d}, nil
	}

	descPremiumDecimal, err := convertToDecimal(desc_premium)
	if err != nil {
		http.Error(w, "invalid description premium!", http.StatusBadRequest)
		return
	}
	descDiscountDecimal, err := convertToDecimal(desc_discount)
	if err != nil {
		http.Error(w, "invalid description discount!", http.StatusBadRequest)
		return
	}
	descAdminCostDecimal, err := convertToDecimal(desc_admin_cost)
	if err != nil {
		http.Error(w, "invalid description admin cost!", http.StatusBadRequest)
		return
	}
	descRiskManagementDecimal, err := convertToDecimal(desc_risk_management)
	if err != nil {
		http.Error(w, "invalid description risk management!", http.StatusBadRequest)
		return
	}
	descBrokageDecimal, err := convertToDecimal(desc_brokage)
	if err != nil {
		http.Error(w, "invalid description brokage!", http.StatusBadRequest)
		return
	}
	descPPHDecimal, err := convertToDecimal(desc_pph)
	if err != nil {
		http.Error(w, "invalid description pph!", http.StatusBadRequest)
		return
	}

	// insAmountDecimal, err := convertToDecimal(ins_amount)
	// if err != nil {
	// 	http.Error(w, "invalid insured amount!", http.StatusBadRequest)
	// 	return
	// }

	// sumInsAmountDecimal, err := convertToDecimal(sum_ins_amount)
	// if err != nil {
	// 	http.Error(w, "invalid insured amount!", http.StatusBadRequest)
	// 	return
	// }

	/*
			GET TOKEN USER LOGIN
		=================================
	*/
	userID, err := auth.GetTokenUserLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	invoices_M := models.Invoice{}
	invoice_ID_Generate, err := models.GenerateInvoiceID(server.DB, models.Type(typeInvoice))
	if err != nil {
		http.Error(w, "Create invoices fail - ID Fail!", http.StatusSeeOther)
	}
	invoices := &models.Invoice{
		Invoice_ID:           invoice_ID_Generate,
		UserID:               userID,
		Type:                 models.Type(typeInvoice),
		Recipient:            recipient,
		Address:              address,
		Desc_Premium:         descPremiumDecimal,
		Desc_Discount:        descDiscountDecimal,
		Desc_Admin_Cost:      descAdminCostDecimal,
		Desc_Risk_Management: descRiskManagementDecimal,
		Desc_Brokage:         descBrokageDecimal,
		Desc_PPH:             descPPHDecimal,
		Policy_Number:        policy_number,
		Name_Of_Insured:      name_of_insured,
		Address_Of_Insured:   address_of_insured,
		Type_Of_Insurance:    type_of_insurance,
		Period_Start:         p_Start,
		Period_End:           p_End,
		Terms_Of_Period:      terms_of_period,
		Remarks:              remarks,
	}
	installments := []models.Installment{}
	sum_insured := []models.Sum_Insured_Details{}

	err = json.Unmarshal([]byte(r.FormValue("installments")), &installments)
	if err != nil {
		http.Error(w, "Parsing installments fail!", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(r.FormValue("sum_insured")), &sum_insured)
	if err != nil {
		http.Error(w, "Parsing sum insured fail1", http.StatusBadRequest)
		return
	}

	err = invoices_M.UpdateInvoices(server.DB, invoiceID, installments, sum_insured)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "update invoice fail", http.StatusBadRequest)
		return
	}
	fmt.Println(invoices)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoices)

}

func (server *Server) DeletedInvoicesAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoice_ID := vars["invoice_id"]

	invoices := &models.Invoice{}
	err := invoices.DeletedInvoices(server.DB, invoice_ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Delete Fail!", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Invoice deleted successfully"))
}

func (server *Server) downloadInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoice_id"]

	invoiceModel := models.Invoice{}
	userModel := models.User{}
	invoice, err := invoiceModel.GetInvoiceByIDmodel(server.DB, invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Invoice not found!", http.StatusBadRequest)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, userModel.CompanyName)

	// Tambahkan informasi invoice
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Type: "+string(invoice.Type))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Desc_Premium: "+invoice.Desc_Premium.String())
	pdf.Ln(10)
	pdf.Cell(40, 10, "Desc_Discount: "+invoice.Desc_Discount.String())
	pdf.Ln(10)
	pdf.Cell(40, 10, "Desc_Discount: "+invoice.Desc_Admin_Cost.String())
	pdf.Ln(10)
	pdf.Cell(40, 10, "Desc_Discount: "+invoice.Desc_Risk_Management.String())
	pdf.Ln(10)
	pdf.Cell(40, 10, "Desc_Discount: "+invoice.Desc_Brokage.String())
	pdf.Ln(10)
	pdf.Cell(40, 10, "Desc_Discount: "+invoice.Desc_PPH.String())
	pdf.Ln(10)

	total, err := invoiceModel.calculateTotalDesc(invoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Summary fail!", http.StatusBadRequest)
		return
	}

	pdf.Cell(40, 10, "Policy Number: "+invoice.Policy_Number)
	pdf.Ln(10)

	// Tampilkan informasi dari Sum_Insured_Details
	sum_ins := models.Sum_Insured_Details{}
	sumInsuredDetails := sum_ins.GetSumInsByInvoiceID(invoiceID)
	for _, detail := range sumInsuredDetails {
		pdf.Cell(40, 10, "Items Name: "+detail.Items_Name)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Amount : "+detail.Sum_Insured_Amount)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Notes : "+detail.Notes)
		pdf.Ln(10)
		// Tambahkan informasi lainnya dari Sum_Insured_Details
	}

}
