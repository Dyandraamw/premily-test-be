package controllers

import (
	"fmt"
	"net/http"
)

func (server *Server)Home(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("Selamat datang di premily")
}