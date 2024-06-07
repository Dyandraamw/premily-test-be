package controllers

import (
	"net/http"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
)

func (server *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	var user models.User
	_, err := user.FindEmailRegis(server.DB, email)
	if err != nil {
		http.Error(w, "Email not found!"+err.Error(), http.StatusNotFound)
		return
	}

	resetToken := uuid.New().String()
	resetTokenExp := time.Now().Add(1 * time.Hour)

	user.Reset_Token = resetToken
	user.Reset_TokenExp = resetTokenExp

	err = server.DB.Save(user).Error
	if err != nil{
		http.Error(w, "Set token for reset fail!"+err.Error(), http.StatusBadRequest)
		return
	}

	//send reset email
	to := user.Email
}
