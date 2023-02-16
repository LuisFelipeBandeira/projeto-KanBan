package routes

import (
	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/server"
)

func SetRoutes(mux *mux.Router) {
	mux.HandleFunc("/createcard", server.CreateCard).Methods("POST")
	mux.HandleFunc("/cards", server.ListCards).Methods("GET")
	mux.HandleFunc("/cards/{cardid}", server.DeleteCard).Methods("DELETE")
}
