package controllers

import (
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("Selamat datang di premily")
}