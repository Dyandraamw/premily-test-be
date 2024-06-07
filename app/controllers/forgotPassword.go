package controllers

import (
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	subject := "Password Reset Request"
	body := "To reset your password, please click the link below:\n\n" +
			os.Getenv("APP_PORT")+ resetToken
	
	err = sendEmail(to, subject, body)
	if err != nil {
		http.Error(w, "Failed to send reset email: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password reset email sent"))
}


func sendEmail(to, subject, body string) error {
	from := os.Getenv("APP_EMAIL")
	password := os.Getenv("APP_EMAIL_PASSWORD")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return err
}

func (server *Server) ResetPasswordAction(w http.ResponseWriter, r *http.Request){
	token := r.FormValue("token")
	newPassword := r.FormValue("new_password")

	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	if newPassword == "" {
		http.Error(w, "New password is required", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := server.DB.First(&user, "reset_token = ?", token).Error; err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	if user.Reset_TokenExp.Before(time.Now()){
		http.Error(w, "Token has expired", http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(newPassword); err != nil {
		http.Error(w, "Please make the secure password"+err.Error(), http.StatusBadRequest)
		return
	}

	makeNewPassword, err :=bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(makeNewPassword)
	user.Reset_Token = ""
	user.Reset_TokenExp = time.Time{}

	err = server.DB.Save(user).Error
	if err != nil{
		http.Error(w, "Failed new password!"+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password has been reset successfully"))
}