package controllers

import (
	// "github.com/frangklynndruru/premily_backend/app/controllers"
	"github.com/gorilla/mux"
)



func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()
	// server.Router.HandleFunc("/", server.Home).Methods("GET")
	server.Router.HandleFunc("/invoice-list", server.Invoice).Methods("GET")

	// server.Router.HandleFunc("/login", server.LoginPage).Methods("GET")
	server.Router.HandleFunc("/login", server.LoginAction).Methods("POST")
}
