package userDB

import (
	"database/sql"
	"log"
)

func CreateUser(db *sql.DB, username string) (int, error) {
	statement, err := db.Prepare("INSERT INTO User(username) VALUES (?)")
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer statement.Close()

	result, err := statement.Exec(username)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	id := int(id64)

	return id, nil
}
