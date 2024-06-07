package controllers

import(
	"net/http"
	"github.com/frangklynndruru/premily_backend/app/models"
)

func (server *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	var user models.User
	_, err := user.FindEmailRegis(server.DB, email)
	if err != nil {
		http.Error(w, "Email not found!"+err.Error(), http.StatusNotFound)
		return
	}

	
}
