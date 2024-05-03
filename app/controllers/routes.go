package controllers

import (
	// "github.com/frangklynndruru/premily_backend/app/controllers"
	"github.com/gorilla/mux"
)

// import "github.com/frangklynndruru/premily_backend/app/controllers"

// func (server *Server) initializeRoutes() {
// 	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
// }

func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()
	// server.Router.HandleFunc("/", server.Home).Methods("GET")
}
