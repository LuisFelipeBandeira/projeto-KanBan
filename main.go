package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/configuration"
	"github.com/projeto-BackEnd/routes"
)

func main() {
	configuration.ConnectDb()

	router := mux.NewRouter()

	routes.SetRoutes(router)

	fmt.Println("Listening in port :8080")
	http.ListenAndServe(":8080", router)
}
