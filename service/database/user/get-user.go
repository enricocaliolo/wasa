package user

import (
	"database/sql"
	"log"
)

func GetUser(db *sql.DB, username string) int {
	statement, err := db.Prepare("SELECT user_id from User WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	var id int
	err = statement.QueryRow(username).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}
