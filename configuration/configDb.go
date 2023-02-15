package configuration

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	var db *sql.DB
	var errConnect error

	db, errConnect = sql.Open("MySql", "root:94647177_Mc@tpc(localhost:33006)/functionarys")
	if errConnect != nil {
		log.Fatal("Error to connect to database: ", errConnect.Error())
	}

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal("Error to Ping in database: ", errPing.Error())
		return nil, errPing
	}

	return db, nil
}
