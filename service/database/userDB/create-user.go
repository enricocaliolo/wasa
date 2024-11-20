package userDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func CreateUser(db *sql.DB, username string) models.User {
	sqlStatement := "INSERT INTO User(username) VALUES (?)"
	result, err := db.Exec(sqlStatement, username)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	err = db.QueryRow("SELECT user_id, username, icon, created_at FROM User WHERE user_id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.Icon, // This will populate the icon from the database
		&user.Created_at,
	)
	if err != nil {
		return models.User{}
	}

	return user
}
