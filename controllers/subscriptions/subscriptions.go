package subscriptions

import (
	"encoding/json"
	"net/http"
	"sensibull-test/services"
	"sensibull-test/structures/subscriptions"

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
	w.Header().Set("Content-Type", "application/json")

	var subscriptionService services.SubscriptionService
	resp, err := subscriptionService.GetByUserName(userName)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, err.Error(), http.StatusNotFound)
}

func GetByUserNameAndDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	date := vars["date"]
	w.Header().Set("Content-Type", "application/json")

	var subscriptionService services.SubscriptionService
	resp, err := subscriptionService.GetByUserNameAndDate(userName, date)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, err.Error(), http.StatusNotFound)
}
