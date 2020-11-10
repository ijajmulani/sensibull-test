package users

import (
	"encoding/json"
	"net/http"
	"strings"

	"sensibull-test/services"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userName = strings.Trim(userName, " ")
	w.Header().Set("Content-Type", "application/json")
	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userService services.UserService
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
	w.Header().Set("Content-Type", "application/json")
	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var userService services.UserService
	if err := userService.Add(userName); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
}
