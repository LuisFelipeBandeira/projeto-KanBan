package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/projeto-BackEnd/configuration"
	"github.com/projeto-BackEnd/model"
	"github.com/projeto-BackEnd/security"
)

func CreateCard(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to get requisition's body"))
		return
	}

	var card *model.Card

	errUnmarshal := json.Unmarshal(body, &card)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do Unmarshal"))
		return
	}

	if card.Company != "Q2BANK" && card.Company != "Q2PAY" && card.Company != "Q2INGRESSOS" {
		card.Company = "Q2PAY"
	}

	Db, errConnect := configuration.ConnectDb()
	if errConnect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to ping in database"))
		return
	}

	defer Db.Close()

	statement, errPrepare := Db.Prepare("Insert into cards (Company, Description, DateLimit, HourLimit) Values (?, ?, ?, ?)")
	if errPrepare != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do prepare of insert"))
		return
	}

	defer statement.Close()

	resultInsert, errInsert := statement.Exec(card.Company, card.Desc, card.DateLimit, card.HourLimit)
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
	w.WriteHeader(http.StatusCreated)
}

func ListCards(w http.ResponseWriter, r *http.Request) {
	db, err := configuration.ConnectDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to ping in database"))
		return
	}

	defer db.Close()

	resultSelect, errSelect := db.Query("SELECT Id, Company, Description, DateLimit, HourLimit FROM cards")
	if errSelect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do query"))
		return
	}

	defer resultSelect.Close()

	var cards []model.Card

	for resultSelect.Next() {
		var card model.Card

		if erroScan := resultSelect.Scan(&card.Id, &card.Company, &card.Desc, &card.DateLimit, &card.HourLimit); erroScan != nil {
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
		if erroScan := resultSelect.Scan(&card.Id, &card.Company, &card.Desc, &card.DateLimit, &card.HourLimit); erroScan != nil {
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

func FinishCard(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	ID, errConvertId := strconv.ParseInt(variables["cardid"], 10, 32)
	if errConvertId != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error to convert ID"))
		return
	}

	DB, errToConnectDatabase := configuration.ConnectDb()
	if errToConnectDatabase != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error connecting to DataBase"))
		return
	}

	defer DB.Close()

	var qtdLine int

	DB.QueryRow("SELECT Count(*) FROM cards WHERE ID = ?", ID).Scan(&qtdLine)

	if qtdLine == 1 {
		_, errUpdate := DB.Query("UPDATE cards SET finished = 1 WHERE Id = ?", ID)
		if errUpdate != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error to do UPDATE"))
			return
		}

		w.WriteHeader(200)
		return

	} else {
		w.WriteHeader(404)
		w.Write([]byte("Card not found"))
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
		w.Write([]byte(errPrepare.Error()))
		return
	}

	DB, errConectDB := configuration.ConnectDb()
	if errConectDB != nil {
		w.WriteHeader(500)
		w.Write([]byte(errConectDB.Error()))
		return
	}

	defer DB.Close()

	var resultCount int

	DB.QueryRow("SELECT Count(*) FROM users WHERE username = ?", user.Username).Scan(&resultCount)

	if resultCount != 0 {
		w.WriteHeader(400)
		w.Write([]byte("Usuário já existe."))
		return
	} else {
		passwordHashed, errHash := security.HashPassword(user.Password)
		if errHash != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error to hash password: " + errHash.Error()))
		}

		user.Password = string(passwordHashed)

		prepare, errPrepareInsert := DB.Prepare("Insert INTO users (name, username, password) Values (?, ?, ?)")
		if errPrepareInsert != nil {
			w.WriteHeader(500)
			w.Write([]byte(errPrepareInsert.Error()))
			return
		}

		defer prepare.Close()

		_, errExecPrepare := prepare.Exec(user.Nome, user.Username, user.Password)
		if errExecPrepare != nil {
			w.WriteHeader(500)
			w.Write([]byte(errExecPrepare.Error()))
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
		w.Write([]byte(errGetUser.Error()))
		return
	}

	var user *model.User

	errUnmarshal := json.Unmarshal(body, &user)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to do Unmarshal"))
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

	db.QueryRow("SELECT Count(*) From users Where username = ?", user.Username).Scan(&resultCount)

	if resultCount != 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	} else {
		var userDB model.User

		db.QueryRow("SELECT Id, password FROM users WHERE username = ?", user.Username).Scan(&userDB.Id, &userDB.Password)

		errComparePassword := security.VerificationPasswordAndHash(user.Password, userDB.Password)
		if errComparePassword != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Username or password's invalid"))
			return
		}

		tokenLogin, errToGenerateToken := security.GenerateJsonWebToken(userDB.Id)
		if errToGenerateToken != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error to do login: Error to generate Json Web Token"))
			w.Write([]byte(errToGenerateToken.Error()))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(tokenLogin)
		return
	}

}
