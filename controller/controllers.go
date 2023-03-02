package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/model"
	"github.com/projeto-BackEnd/repository"
	"github.com/projeto-BackEnd/security"
)

func CreateCard(w http.ResponseWriter, r *http.Request) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to get requisition's body")
		return
	}

	var card *model.Card

	errUnmarshal := json.Unmarshal(body, &card)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to do Unmarshal")
		return
	}

	if card.Company != "Q2BANK" && card.Company != "Q2PAY" && card.Company != "Q2INGRESSOS" {
		card.Company = "Q2PAY"
	}

	err := repository.InsertCard(card.Company, card.Desc, card.DateLimit, card.HourLimit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("Task inserida")
	w.WriteHeader(200)
}

func ListCards(w http.ResponseWriter, r *http.Request) {

	resultSelect, errSelect := repository.ListCards()
	if errSelect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to do query")
		return
	}

	var cards []model.Card

	for resultSelect.Next() {
		var card model.Card

		if erroScan := resultSelect.Scan(&card.Id, &card.Company, &card.Desc, &card.DateLimit, &card.HourLimit); erroScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Error to do Scan in database")
			return
		}

		cards = append(cards, card)
	}

	w.WriteHeader(200)

	errEncoder := json.NewEncoder(w).Encode(cards)
	if errEncoder != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to return json")
		return
	}
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, errGetId := strconv.ParseInt(params["cardid"], 10, 32)
	if errGetId != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error to get ID")
		fmt.Println("Error to get ID: ", errGetId.Error())
		return
	}

	resultCount, errConnect := repository.CardExist(int(ID))
	if errConnect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to ping in database")
		return
	}

	if resultCount == 1 {
		errDelete := repository.DeleteCard(int(ID))
		if errDelete != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errDelete.Error())
			return
		}

		w.WriteHeader(200)
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("Card not found")
		return
	}
}

func ListCardUsingId(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	ID, err := strconv.ParseInt(variables["cardid"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to GET Id")
		return
	}

	var card model.Card

	resultSelect, errSelect := repository.ListOneCard(int(ID))
	if errSelect != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error to do Select")
		return
	}

	for resultSelect.Next() {
		if erroScan := resultSelect.Scan(&card.Id, &card.Company, &card.Desc, &card.DateLimit, &card.HourLimit); erroScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Error to do Scan of select's result")
			return
		}
	}

	w.WriteHeader(200)
	errEncoder := json.NewEncoder(w).Encode(card)
	if errEncoder != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to return json")
		return
	}
}

func FinishCard(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	ID, errConvertId := strconv.ParseInt(variables["cardid"], 10, 32)
	if errConvertId != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Error to convert ID")
		return
	}

	qtdLine, err := repository.CardExist(int(ID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if qtdLine == 1 {
		errUpdate := repository.FinishCard(int(ID))
		if errUpdate != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errUpdate.Error())
			return
		}
	} else {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode("Card not found")
		return
	}
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	body, errGetBody := io.ReadAll(r.Body)
	if errGetBody != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user *model.User

	errUnmarshal := json.Unmarshal(body, &user)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errPrepare := user.Prepare()
	if errPrepare != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errPrepare.Error())
		return
	}

	resultCount, errValidUser := repository.UserExist(user.Username)
	if errValidUser != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errValidUser.Error())
		return
	}

	if resultCount != 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Usuário existente.")
		return
	} else {
		errInsert := repository.InsertUser(user.Nome, user.Username, user.Password)
		if errInsert != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errInsert.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, errGetUser := io.ReadAll(r.Body)
	if errGetUser != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errGetUser.Error())
		return
	}

	var user *model.User

	errUnmarshal := json.Unmarshal(body, &user)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error to do Unmarshal")
		return
	}

	resultCount, errVerificUser := repository.UserExist(user.Username)
	if errVerificUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errVerificUser.Error())
		return
	}

	if resultCount != 1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User not found")
		return
	} else {
		var userDB model.User

		resultSelect, errLogin := repository.Login(user.Username)
		if errLogin != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errLogin.Error())
			return
		}

		resultSelect.Scan(&userDB.Id, &userDB.Password)

		errComparePassword := security.VerificationPasswordAndHash(user.Password, userDB.Password)
		if errComparePassword != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Username or password's invalid")
			return
		}

		tokenLogin, errToGenerateToken := security.GenerateJsonWebToken(userDB.Id)
		if errToGenerateToken != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode("Error to do login: Error to generate Json Web Token")
			json.NewEncoder(w).Encode(errToGenerateToken.Error())
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(tokenLogin)
		return
	}

}
