package main

import (
	"net/http"
	"os"
	"sensibull-test/routers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := routers.SetupRouter()
	http.ListenAndServe(":"+os.Getenv("HTTPPORT"), r)
}
