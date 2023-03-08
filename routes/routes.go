package routes

import (
	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/controller"
	"github.com/projeto-BackEnd/middlewares"
)

func SetRoutes(mux *mux.Router) {
	mux.Use(middlewares.JsonMiddleware)
	mux.HandleFunc("/createcard", controller.CreateCard).Methods("POST")
	mux.HandleFunc("/cards", controller.ListCards).Methods("GET")
	mux.HandleFunc("/cards/{cardid}", controller.DeleteCard).Methods("DELETE")
	mux.HandleFunc("/card/{cardid}", controller.ListCardUsingId).Methods("GET")
	mux.HandleFunc("/card/{cardid}", controller.FinishCard).Methods("PUT")
	mux.HandleFunc("/createuser", controller.InsertUser).Methods("POST")
	mux.HandleFunc("/login", controller.Login).Methods("POST")

}
