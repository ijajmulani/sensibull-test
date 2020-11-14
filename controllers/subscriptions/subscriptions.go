package subscriptions

import (
	"encoding/json"
	"net/http"
	"sensibull-test/services"
	"sensibull-test/structures/subscriptions"
	"strings"

	"github.com/gorilla/mux"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var args subscriptions.PostArgs

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil || args.UserName == "" || args.StartDate == "" || args.PlanName == "" {
		http.Error(w, "Request not well formed", http.StatusBadRequest)
		return
	}

	var subscriptionService services.SubscriptionService
	if err := subscriptionService.Post(args); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetByUserName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userName = strings.Trim(userName, " ")
	w.Header().Set("Content-Type", "application/json")
	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var subscriptionService services.SubscriptionService
	if resp, err := subscriptionService.GetByUserName(userName); err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func GetByUserNameAndDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	date := vars["date"]
	userName = strings.Trim(userName, " ")
	w.Header().Set("Content-Type", "application/json")
	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var subscriptionService services.SubscriptionService
	if resp, err := subscriptionService.GetByUserNameAndDate(userName, date); err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
