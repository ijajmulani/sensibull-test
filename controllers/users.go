package controllers

import (
	"encoding/json"
	"net/http"
	"sensibull-test/services"
)

func Get(w http.ResponseWriter, r *http.Request) {
	var userService services.UserService
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := userService.List()
	json.NewEncoder(w).Encode(resp)
}

func Add(w http.ResponseWriter, r *http.Request) {
	var userService services.UserService
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := userService.Add()
	json.NewEncoder(w).Encode(resp)
}
