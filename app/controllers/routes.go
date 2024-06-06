package controllers

import (
	"github.com/frangklynndruru/premily_backend/app/controllers/middlewares"
	"github.com/gorilla/mux"
)

func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()
	// server.Router.HandleFunc("/login", server.LoginPage).Methods("GET")

	/*
		Authentication  Sign-In & Sign-Up
	*/
	server.Router.HandleFunc("/sign-in", server.SignInAction).Methods("POST")
	server.Router.HandleFunc("/sign-up", server.SignUpAction).Methods("POST")

	api := server.Router.PathPrefix("/api").Subrouter()
	api.Use(middlewares.JwtMiddleware)

	api.HandleFunc("/sign-out", server.SignOutAction).Methods("POST")

	/* === HANDLE ROUTES FOR ACCESS CONTROL ACTION === */
	api.HandleFunc("/unverified-users", server.GetUnverifiedUserAction).Methods("GET")
	api.HandleFunc("/unrole-users", server.GetUnroleUserAction).Methods("GET")
	api.HandleFunc("/verify-user/{user_id}", server.VerifyUserAction).Methods("POST")
	api.HandleFunc("/set-role/{user_id}", server.SetUserRoleAction).Methods("POST")


	api.HandleFunc("/user/{user_id}", server.GetUserAction).Methods("GET")
	api.HandleFunc("/update-profile-picture/{user_id}", server.UpdatePhotoProfile).Methods("PUT")


	/* === HANDLE ROUTES FOR INVOICE === */
	api.HandleFunc("/invoice-list", server.Invoice).Methods("GET")
	api.HandleFunc("/invoice/{invoice_id}", server.GetInvoiceByID).Methods("GET")
	api.HandleFunc("/create-invoices", server.CreateInvoicesAction).Methods("POST")
	api.HandleFunc("/update-invoices/{invoice_id}", server.UpdateInvoices).Methods("PUT")
	api.HandleFunc("/delete-invoices/{invoice_id}", server.DeletedInvoicesAction).Methods("DELETE")
	
	
	/* === HANDLE ROUTES FOR STATEMENT OF ACCOUNTS ( S O A ) === */
	api.HandleFunc("/retrive-soa", server.GetSoaResponseList).Methods("GET")
	api.HandleFunc("/get-items-list/{soa_id}", server.GetItemsBySoaID).Methods("GET")
	api.HandleFunc("/create-soa", server.CreateSoaAction).Methods("POST")
	api.HandleFunc("/add-items/{soa_id}", server.AddItemSoaAction).Methods("POST")
	api.HandleFunc("/edit-items/{soa_details_id}", server.UpdateItemSoaAction).Methods("PUT")
	api.HandleFunc("/delete-soa/{soa_id}", server.DeleteSoaAction).Methods("DELETE")


	/* === HANDLE ROUTES FOR ADJUSTMENT === */
	api.HandleFunc("/add-adjustment/{payment_status_id}", server.AddAjustment).Methods("POST")


	/* === HANDLE ROUTES FOR PAYMENT STATUS === */
	api.HandleFunc("/payment-data/{payment_status_id}", server.GetPaymentData).Methods("GET")
	api.HandleFunc("/create-new-payment-status", server.CreatePaymentStatus).Methods("POST")
	api.HandleFunc("/add-payment", server.AddPayment).Methods("POST")
	api.HandleFunc("/edit-payment/{pay_detail_id}", server.UpdatePaymentDetails).Methods("PUT")

}
