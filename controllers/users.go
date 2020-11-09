package controllers

import (
	"encoding/json"
	"net/http"
	"sensibull-test/services"
	"strings"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userName = strings.Trim(userName, " ")
	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userService services.UserService
	w.Header().Set("Content-Type", "application/json")
	if resp, err := userService.Get(userName); err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userName = strings.Trim(userName, " ")
	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	var userService services.UserService
	w.Header().Set("Content-Type", "application/json")
	if err := userService.Add(userName); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
}
