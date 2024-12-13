package userDB

import (
	"database/sql"
)

func GetUser(db *sql.DB, username string) (string, error) {
	statement, err := db.Prepare("SELECT user_id from User WHERE username = ?")
	if err != nil {
		return"", err
	}
	defer statement.Close()
	var id int
	err = statement.QueryRow(username).Scan(&id)
	if err != nil {
		return"", err
	}

	return username, nil
}
