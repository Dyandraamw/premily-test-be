package controllers

import (
	// "github.com/frangklynndruru/premily_backend/app/controllers"
	"github.com/gorilla/mux"
)



func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()
	// server.Router.HandleFunc("/", server.Home).Methods("GET")
	server.Router.HandleFunc("/invoice-list", server.Invoice).Methods("GET")

	/*
		Authentication  Sign-In & Sign-Up
	*/
	// server.Router.HandleFunc("/login", server.LoginPage).Methods("GET")
	server.Router.HandleFunc("/sign-in", server.SignInAction).Methods("POST")
	server.Router.HandleFunc("/sign-out", server.SignOutAction).Methods("GET")
	server.Router.HandleFunc("/sign-up", server.SignUpAction).Methods("POST")

}
