package userDB

import (
	"database/sql"
)

func CreateUser(db *sql.DB, username string) (int, error) {
	statement, err := db.Prepare("INSERT INTO User(username) VALUES (?)")
	if err != nil {
	return -1, err
	}
	defer statement.Close()

	result, err := statement.Exec(username)
	if err != nil {
	return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
	return -1, err
	}

	return int(id), nil

}
