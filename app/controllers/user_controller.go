package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	
	"golang.org/x/crypto/bcrypt"
)
// func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request){
// 	http.Redirect(w, r, "/login", http.StatusOK)

// }
func (server *Server) LoginAction(w http.ResponseWriter, r *http.Request){
	userModel := models.User{}
	
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := userModel.FindByEmail(server.DB, email)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 
	}

	var compare_password = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if compare_password == nil {
		// store.Set("email", user.Email)
		// http.Redirect(w, r, "/login", http.StatusSeeOther )
		return 
	}
	session, _ := store.Get(r, sessionUser)
	session.Values["user_id"] = user.UserID
	session.Save(r, w)

	// http.Redirect(w, r, "/", http. StatusSeeOther)

	data, _:= json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}