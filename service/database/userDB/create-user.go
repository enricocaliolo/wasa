package userDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func CreateUser(db *sql.DB, username string) models.User {
	statement, err := db.Prepare("INSERT INTO User(username) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	result, err := statement.Exec(username)
	if err != nil {
		log.Fatal(err)
	}

	id64, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	id := int(id64)

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
