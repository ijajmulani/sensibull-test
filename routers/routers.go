package routers

import (
	"net/http"

	apiController "sensibull-test/controllers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/add", apiController.Add).Methods(http.MethodGet)
	r.HandleFunc("/user", apiController.Get).Methods(http.MethodGet)
	return r
}
