package user

import (
	"database/sql"
	"log"
)

func CreateUser(db *sql.DB, username string) int {
	sqlStatement := "INSERT INTO User(username) VALUES (?)"
	result, err := db.Exec(sqlStatement, username)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return int(id)
}
