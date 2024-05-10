package controllers

import (
	"encoding/json"
	_ "fmt"
	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	"golang.org/x/crypto/bcrypt"
)

// func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request){
// 	http.Redirect(w, r, "/login", http.StatusOK)
//render frontend

// }
func (server *Server) SignInAction(w http.ResponseWriter, r *http.Request){
	userModel := models.User{}
	
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := userModel.FindByEmail(server.DB, email, password)
	if err != nil{
		
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 
	}
	// fmt.Println("Lolos Email")
	
	// if !verifyPassword(password, user.Password) {
		
	// 	http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
	// 	// http.Redirect(w, r, "/login", http.StatusSeeOther )
	// 	return 
	// }
	var compare_password = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if compare_password == nil {
		
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		// http.Redirect(w, r, "/login", http.StatusSeeOther )
		return 
	}
	// fmt.Println("Lolos PW")
	session, _ := store.Get(r, sessionUser)
	session.Values["user_id"] = user.UserID
	session.Save(r, w)

	// http.Redirect(w, r, "/", http. StatusSeeOther)

	data, _:= json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (server *Server) SignOutAction(w http.ResponseWriter, r *http.Request){

	session, _ := store.Get(r, sessionUser)

	session.Values["user_id"]= nil
	session.Save(r,w)

	// http.Redirect(w, r, "/", http.StatusOK)
	
	
}

func (sever *Server) SignUpAction(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	email := r.FormValue("email")
	phoneNumber := r.FormValue("phoneNumber")
	password := r.FormValue("password")
	company := r.FormValue("company")
}