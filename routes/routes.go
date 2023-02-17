package routes

import (
	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/controller"
)

func SetRoutes(mux *mux.Router) {
	mux.HandleFunc("/createcard", controller.CreateCard).Methods("POST")
	mux.HandleFunc("/cards", controller.ListCards).Methods("GET")
	mux.HandleFunc("/cards/{cardid}", controller.DeleteCard).Methods("DELETE")
	mux.HandleFunc("/card/{cardid}", controller.ListCardUsingId).Methods("GET")
	mux.HandleFunc("/card/{cardid}", controller.FinishCard).Methods("PUT")
}
