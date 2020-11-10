package routers

import (
	"net/http"

	subscriptionsController "sensibull-test/controllers/subscriptions"
	usersController "sensibull-test/controllers/users"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/{userName}", usersController.Put).Methods(http.MethodPut)
	r.HandleFunc("/user/{userName}", usersController.Get).Methods(http.MethodGet)

	// r.HandleFunc("/subscription/", usersController.Put).Methods(http.MethodPut)
	r.HandleFunc("/subscription/{userName}", subscriptionsController.List).Methods(http.MethodGet)
	r.HandleFunc("/subscription/{userName}/{date}", subscriptionsController.GetWithDate).Methods(http.MethodGet)
	r.HandleFunc("/subscription", subscriptionsController.Post).Methods(http.MethodPost)
	return r
}
