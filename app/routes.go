package app

import "github.com/frangklynndruru/premily_backend/app/controllers"

// import "github.com/frangklynndruru/premily_backend/app/controllers"

// func (server *Server) initializeRoutes() {
// 	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
// }

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
