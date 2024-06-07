package controllers

import (
	"encoding/json"
	"time"

	"net/http"
	"os"
	"path/filepath"

	"github.com/frangklynndruru/premily_backend/app/models"
	"github.com/gorilla/mux"
)

func (server *Server) UpdatePhotoProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]

	var existPhoto models.User
	if err := server.DB.First(&existPhoto, "user_id = ?", user_id).Error; err != nil {
		http.Error(w, "Profile not found! "+err.Error(), http.StatusNotFound)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a directory to save the uploaded files
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Create a new file in the uploads directory
	filePath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the new file
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the user profile with the new image path
	existPhoto.Image = filePath
	existPhoto.Updated_At = time.Now()

	err = existPhoto.UpdateUserPicture(server.DB, user_id)
	if err != nil {
		http.Error(w, "Failed to update profile picture: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal the updated user profile into JSON
	data, err := json.Marshal(existPhoto)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
