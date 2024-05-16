package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"regexp"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request){
// 	http.Redirect(w, r, "/login", http.StatusOK)
//render frontend

// }
func (server *Server) GetUserAction(w http.ResponseWriter, r *http.Request){
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

	userDetails := struct{
		Username    string       `json:"username"`
        Name        string       `json:"name"`
        Email       string       `json:"email"`
        Phone       string       `json:"phone"`
        CompanyName string       `json:"company_name"`
        Role        models.Role  `json:"role"`
	}{
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
	if user.Role != models.StaffRole || user.Role != models.AdminRole || user.Role != models.AccessControlRole {
        http.Error(w, "Access denied", http.StatusForbidden)
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

	if err := ValidatePassword(password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		Role: "pending",
		Verified: false,
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
func (server *Server) GetUnverifiedUserAction(w http.ResponseWriter, r *http.Request){
	userModel := models.User{}
	users, err := userModel.GetUnverifiedUser(server.DB)
	if err != nil{
		http.Error(w, "Failed to get unverified users", http.StatusBadRequest)
	}

	data, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (server *Server) VerifyAndSetRoleUserAction(w http.ResponseWriter, r *http.Request){
	user_id := r.FormValue("user_id")
	role := r.FormValue("role")

	if user_id == "" || role == "" {
		http.Error(w, "Please fill ID and set the role of user", http.StatusBadRequest)
		return
	}

	userModel := models.User{}
	if err := userModel.VerifyAndSetUserRole(server.DB, user_id, models.Role(role)); err != nil {
		http.Error(w, "Verify user fail!", http.StatusConflict)
	}
	w.WriteHeader(http.StatusOK)
    w.Write([]byte("User successfully verified"))
}

// func ValidatePassword(password string) error {
// 	if len(password) < 8 {
// 		return fmt.Errorf("Password must be at least 8 characters!")
// 	}

// 	// Regular expression to check for at least one lowercase letter, one uppercase letter, one number, and one special character.
// 	match, err := regexp.MatchString(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]+$`, password)
// 	if err != nil {
// 		return fmt.Errorf("Error while validating password: %v", err)
// 	}

// 	if !match {
// 		return fmt.Errorf("Password must contain uppercase letter, lowercase letter, number, and special character")
// 	}

// 	return nil
// }
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
    }

    // Using simpler regex supported by Go's regexp package
    match := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]+$`)
    if !match.MatchString(password) {
        return fmt.Errorf("password contains invalid characters")
    }

    // Manual checks for required character types
    if !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return fmt.Errorf("password must contain at least one lowercase letter")
    }
    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return fmt.Errorf("password must contain at least one uppercase letter")
    }
    if !regexp.MustCompile(`\d`).MatchString(password) {
        return fmt.Errorf("password must contain at least one digit")
    }
    if !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) {
        return fmt.Errorf("password must contain at least one special character")
    }

    return nil
}