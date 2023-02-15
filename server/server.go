package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/projeto-BackEnd/configuration"
	"github.com/projeto-BackEnd/model"
)

func CreateCard(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to get requisition's body"))
		return
	}

	var card model.Card

	errUnmarshal := json.Unmarshal(body, &card)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do Unmarshal"))
		return
	}

	Db, errConnect := configuration.ConnectDb()
	if errConnect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to ping in database"))
		return
	}

	defer Db.Close()

	statement, errPrepare := Db.Prepare("Insert into cards (title, description) Values (?, ?)")
	if errPrepare != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do prepare"))
		return
	}

	statement.Exec(card.Title, card.Desc)
}
