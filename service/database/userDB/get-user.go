package userDB

import (
	"database/sql"
	"log"
	. "wasa/service/shared/models"
)

func GetUser(db *sql.DB, username string) User {
	statement, err := db.Prepare("SELECT * from User WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	var user User
	err = statement.QueryRow(username).Scan(&user.ID, &user.Username, &user.Icon, &user.Created_at)
	if err != nil {
		return User{
			ID: -1,
		}
	}

	return user
}
