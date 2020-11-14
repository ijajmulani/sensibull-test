package routers

import (
	"net/http"

	subscriptionsController "sensibull-test/controllers/subscriptions"
	usersController "sensibull-test/controllers/users"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/{userName:[a-zA-Z0-9]+}", usersController.Put).Methods(http.MethodPut)
	r.HandleFunc("/user/{userName:[a-zA-Z0-9]+}", usersController.Get).Methods(http.MethodGet)

	// r.HandleFunc("/subscription/", usersController.Put).Methods(http.MethodPut)
	r.HandleFunc("/subscription/{userName:[a-zA-Z0-9]+}", subscriptionsController.GetByUserName).Methods(http.MethodGet)
	r.HandleFunc("/subscription/{userName:[a-zA-Z0-9]+}/{date}", subscriptionsController.GetByUserNameAndDate).Methods(http.MethodGet)
	r.HandleFunc("/subscription", subscriptionsController.Post).Methods(http.MethodPost)
	return r
}
