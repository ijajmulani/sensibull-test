package subscriptions

import (
	"encoding/json"
	"net/http"
	"sensibull-test/services"
	"sensibull-test/structures/subscriptions"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var args subscriptions.PostArgs

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil || args.UserName == "" || args.StartDate == "" || args.PlanID == "" {
		http.Error(w, "Request not well formed", http.StatusBadRequest)
		return
	}

	var subscriptionService services.SubscriptionService
	if err := subscriptionService.Post(args); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
}
