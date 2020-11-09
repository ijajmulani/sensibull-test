package routers

import (
	"net/http"

	apiController "sensibull-test/controllers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/{userName}", apiController.Put).Methods(http.MethodPut)
	r.HandleFunc("/user/{userName}", apiController.Get).Methods(http.MethodGet)
	return r
}
