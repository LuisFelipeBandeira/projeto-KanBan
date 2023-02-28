package repository

import (
	"github.com/projeto-BackEnd/configuration"
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

func UserExist(id int) (int, error) {
	db, errConnectDB := configuration.ConnectDb()
	if errConnectDB != nil {
		return 0, errConnectDB
	}

	defer db.Close()

	var resultCount int

	db.QueryRow("SELECT COUNT(*) FROM cards WHERE id = ?", id).Scan(&resultCount)

	return resultCount, nil
}
