package userDB

import (
	"database/sql"
)

func CreateUser(db *sql.DB, username string) (string, error) {
	statement, err := db.Prepare("INSERT INTO User(username) VALUES (?)")
	if err != nil {
		return "", err
	}
	defer statement.Close()

	_, err = statement.Exec(username)
	if err != nil {
		return "", err
	}

	return username, nil
}
