package controllers
import (
	"encoding/json"
	"net/http"
	
	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"

)


func (server *Server) GetUnverifiedUserAction(w http.ResponseWriter, r *http.Request) {
	userModel := models.User{}
	users, err := userModel.GetUnverifiedUser(server.DB)
	if err != nil {
		http.Error(w, "Failed to get unverified users", http.StatusBadRequest)
	}

	data, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func (server *Server) GetUnroleUserAction(w http.ResponseWriter, r *http.Request) {
	userModel := models.User{}
	users, err := userModel.GetUnroleUser(server.DB)
	if err != nil {
		http.Error(w, "Failed to get unverified users", http.StatusBadRequest)
	}

	data, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (server *Server) VerifyUserAction(w http.ResponseWriter, r *http.Request) {
	//jangan lupa ganti ke vars := mux.vars(r); userID := vars["user_id"]
	vars := mux.Vars(r)
	user_id := vars["user_id"]
	verify := r.FormValue("verify")

	if verify == "" {
		http.Error(w, "Please set the verify of user", http.StatusBadRequest)
		return
	}

	userModel := models.User{}
	if err := userModel.VerifyUser(server.DB, user_id, models.Verify(verify)); err != nil {
		http.Error(w, "Verify user fail!", http.StatusConflict)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully verified"))
}
func (server *Server) SetUserRoleAction(w http.ResponseWriter, r *http.Request) {
	//jangan lupa ganti ke vars := mux.vars(r); userID := vars["user_id"]
	vars := mux.Vars(r)
	user_id := vars["user_id"]
	role := r.FormValue("role")

	if role == "" {
		http.Error(w, "Please fill the role of user", http.StatusBadRequest)
		return
	}

	userModel := models.User{}
	if err := userModel.SetUserRole(server.DB, user_id, models.Role(role)); err != nil {
		http.Error(w, "Verify user fail!", http.StatusConflict)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully set the role"))
}