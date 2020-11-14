package users

import (
	"encoding/json"
	"net/http"
	"strings"

	"sensibull-test/constants"
	"sensibull-test/services"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userName = strings.Trim(userName, " ")
	w.Header().Set("Content-Type", "application/json")

	var userService services.UserService
	if resp, err := userService.Get(userName); err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, constants.UserNotFound, http.StatusNotFound)
}

func Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	w.Header().Set("Content-Type", "application/json")
	var userService services.UserService
	if err := userService.Add(userName); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
}
