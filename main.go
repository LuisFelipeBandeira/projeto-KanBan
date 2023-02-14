package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/configuration"
	"github.com/projeto-BackEnd/controller/routes"
)

func main() {
	configuration.ConnectDb()

	router := mux.NewRouter()

	routes.SetRoutes(router)

	http.ListenAndServe(":8080", router)
}
