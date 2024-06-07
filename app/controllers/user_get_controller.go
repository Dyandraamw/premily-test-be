package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

func (server *Server) GetUserAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	userModel := models.User{}
	user, err := userModel.FindByID(server.DB, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	userDetails := struct {
		Image       string      `json:image`
		Username    string      `json:"username"`
		Name        string      `json:"name"`
		Email       string      `json:"email"`
		Phone       string      `json:"phone"`
		CompanyName string      `json:"company_name"`
		Role        models.Role `json:"role"`
	}{
		Image:       user.Image,
		Username:    user.Username,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		CompanyName: user.CompanyName,
		Role:        user.Role,
	}
	data, _ := json.Marshal(userDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}