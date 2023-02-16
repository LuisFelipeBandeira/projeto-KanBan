package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/configuration"
	"github.com/projeto-BackEnd/model"
)

func CreateCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
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

	statement, errPrepare := Db.Prepare("Insert into cards (Description, DateLimit, HourLimit) Values (?, ?, ?)")
	if errPrepare != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do prepare of insert"))
		return
	}

	defer statement.Close()

	resultInsert, errInsert := statement.Exec(card.Desc, card.DateLimit, card.HourLimit)
	if errInsert != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do insert in database"))
		return
	}

	_, errResult := resultInsert.LastInsertId()
	if errResult != nil {
		fmt.Println(errResult.Error())
	}

	w.Write([]byte("Task inserida"))
	w.WriteHeader(200)
}

func ListCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := configuration.ConnectDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to ping in database"))
		return
	}

	defer db.Close()

	resultSelect, errSelect := db.Query("SELECT * FROM cards")
	if errSelect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do query"))
		return
	}

	defer resultSelect.Close()

	var cards []model.Card

	for resultSelect.Next() {
		var card model.Card

		if erroScan := resultSelect.Scan(&card.Id, &card.Desc, &card.DateLimit, &card.HourLimit); erroScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error to do Scan in database"))
			return
		}

		cards = append(cards, card)
	}

	w.WriteHeader(200)

	errEncoder := json.NewEncoder(w).Encode(cards)
	if errEncoder != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to return json"))
		return
	}
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	ID, errGetId := strconv.ParseInt(params["cardid"], 10, 32)
	if errGetId != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error to get ID"))
		fmt.Println("Error to get ID: ", errGetId.Error())
		return
	}

	db, errConnectDB := configuration.ConnectDb()
	if errConnectDB != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to ping in database"))
		return
	}

	defer db.Close()

	var resultCount int

	db.QueryRow("SELECT COUNT(*) FROM cards WHERE Id = ?", ID).Scan(&resultCount)

	if resultCount == 1 {
		statementDelete, errPrepare := db.Prepare("DELETE FROM cards WHERE Id = ?")
		if errPrepare != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error to do prepare of delete"))
			return
		}

		defer statementDelete.Close()

		_, errDelete := statementDelete.Exec(ID)
		if errDelete != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error to do delete"))
			return
		}

		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
		w.Write([]byte("Task not found"))
	}
}

func ListCardUsingId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	variables := mux.Vars(r)
	ID, err := strconv.ParseInt(variables["cardid"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to GET Id"))
		return
	}

	db, errConnect := configuration.ConnectDb()
	if errConnect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to ping in database"))
		return
	}

	defer db.Close()

	resultSelect, errSelect := db.Query("SELECT * FROM cards WHERE Id = ?", ID)
	if errSelect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do query"))
		return
	}

	defer resultSelect.Close()

	var card model.Card

	for resultSelect.Next() {
		if erroScan := resultSelect.Scan(&card.Id, &card.Desc, &card.DateLimit, &card.HourLimit); erroScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error to do Scan in database"))
			return
		}
	}

	w.WriteHeader(200)

	errEncoder := json.NewEncoder(w).Encode(card)
	if errEncoder != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to return json"))
		return
	}
}
