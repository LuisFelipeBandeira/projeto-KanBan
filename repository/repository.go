package repository

import (
	"database/sql"

	"github.com/projeto-BackEnd/configuration"
	"github.com/projeto-BackEnd/security"
)

func InsertCard(company, description, dateLimit, hourLimit string) error {
	Db, errConnect := configuration.ConnectDb()
	if errConnect != nil {
		return errConnect
	}

	defer Db.Close()

	statement, errPrepare := Db.Prepare("Insert into cards (Company, Description, DateLimit, HourLimit) Values (?, ?, ?, ?)")
	if errPrepare != nil {
		return errPrepare
	}

	defer statement.Close()

	_, errInsert := statement.Exec(company, description, dateLimit, hourLimit)
	if errInsert != nil {
		return errInsert
	}

	return nil
}

func DeleteCard(cardId int) error {
	db, errConnectDB := configuration.ConnectDb()
	if errConnectDB != nil {
		return errConnectDB
	}

	defer db.Close()

	statementDelete, errPrepare := db.Prepare("DELETE FROM cards WHERE Id = ?")
	if errPrepare != nil {
		return errPrepare
	}

	defer statementDelete.Close()

	_, errDelete := statementDelete.Exec(cardId)
	if errDelete != nil {
		return errDelete
	}

	return nil
}

func UserExist(userName string) (int, error) {
	db, errConnectDB := configuration.ConnectDb()
	if errConnectDB != nil {
		return 0, errConnectDB
	}

	defer db.Close()

	var resultCount int

	db.QueryRow("SELECT Count(*) FROM users WHERE username = ?", userName).Scan(&resultCount)

	return resultCount, nil
}

func CardExist(id int) (int, error) {
	db, errConnectDB := configuration.ConnectDb()
	if errConnectDB != nil {
		return 0, errConnectDB
	}

	defer db.Close()

	var resultCount int

	db.QueryRow("SELECT COUNT(*) FROM cards WHERE id = ?", id).Scan(&resultCount)

	return resultCount, nil
}

func ListCards() (*sql.Rows, error) {
	db, err := configuration.ConnectDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	resultSelect, errSelect := db.Query("SELECT Id, Company, Description, DateLimit, HourLimit FROM cards")
	if errSelect != nil {
		return nil, errSelect
	}

	return resultSelect, nil
}

func ListOneCard(id int) (*sql.Rows, error) {
	db, errConnect := configuration.ConnectDb()
	if errConnect != nil {
		return nil, errConnect
	}

	defer db.Close()

	resultSelect, errSelect := db.Query("SELECT * FROM cards WHERE Id = ?", id)
	if errSelect != nil {
		return nil, errSelect
	}

	return resultSelect, nil
}

func FinishCard(id int) error {
	DB, errToConnectDatabase := configuration.ConnectDb()
	if errToConnectDatabase != nil {
		return errToConnectDatabase
	}

	defer DB.Close()

	_, errUpdate := DB.Query("UPDATE cards SET finished = 1 WHERE Id = ?", id)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func InsertUser(nome, username, password string) error {

	db, errConnect := configuration.ConnectDb()
	if errConnect != nil {
		return errConnect
	}

	defer db.Close()

	passwordHashed, errHash := security.HashPassword(password)
	if errHash != nil {
		return errHash
	}

	prepare, errPrepareInsert := db.Prepare("Insert INTO users (name, username, password) Values (?, ?, ?)")
	if errPrepareInsert != nil {
		return errPrepareInsert
	}

	defer prepare.Close()

	_, errExecPrepare := prepare.Exec(nome, username, passwordHashed)
	if errExecPrepare != nil {
		return errExecPrepare
	}

	return nil
}

func Login(username string) (*sql.Row, error) {
	db, errConnectDB := configuration.ConnectDb()
	if errConnectDB != nil {
		return nil, errConnectDB
	}

	defer db.Close()

	result := db.QueryRow("SELECT Id, password FROM users WHERE username = ?", username)

	return result, nil
}
