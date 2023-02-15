package routes

import (
	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/server"
)

func SetRoutes(mux *mux.Router) {
	mux.HandleFunc("/cards", server.CreateCard).Methods("POST")
}
