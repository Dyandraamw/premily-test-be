package controllers

import (
	"encoding/json"
	_ "fmt"
	"net/http"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request){
// 	http.Redirect(w, r, "/login", http.StatusOK)
//render frontend

// }
func (server *Server) SignInAction(w http.ResponseWriter, r *http.Request) {
	userModel := models.User{}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := userModel.FindByEmail(server.DB, email, password)
	if err != nil {

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

	data, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (server *Server) SignOutAction(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, sessionUser)

	session.Values["user_id"] = nil
	session.Save(r, w)

	// http.Redirect(w, r, "/", http.StatusOK)

}

func (server *Server) SignUpAction(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	company := r.FormValue("company")

	if username == "" || name == "" || email == "" || phone == "" || password == "" || confirmPassword == "" || company == "" {
		http.Error(w, "Please fill the required fields!", http.StatusSeeOther)
		return
	}

	userModel := models.User{}
	userRegistered, _ := userModel.FindEmailRegis(server.DB, email)
	if userRegistered != nil {
		http.Error(w, "Email already sign-up!", http.StatusConflict)
		return
	}

	if password != confirmPassword {
		http.Error(w, "Password do not match", http.StatusUnauthorized)
		return
	}
	// hashPassword, _ := MakePassword(password)
	makePassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	params := &models.User{
		UserID:      uuid.New().String(),
		Username:    username,
		Name:        name,
		Email:       email,
		Phone:       phone,
		Password:    string(makePassword),
		CompanyName: company,
	}

	user, err := userModel.CreateUser(server.DB, params)
	if err != nil {
		http.Error(w, "Failed to sign-up!", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
