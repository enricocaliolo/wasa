package userDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func UpdateProfile(db *sql.DB, user models.User) bool {
	if isUsernameTaken(db, user.Username) {
		return false
	}

	statement, err := db.Prepare(`
    UPDATE User
    SET username = COALESCE(?, username),
        icon = COALESCE(?, icon)
    WHERE user_id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	_, err = statement.Exec(user.Username, user.Icon, user.ID)

	return true
}

func isUsernameTaken(db *sql.DB, username string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM User WHERE LOWER(username) = LOWER(?))"
	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return true
	}
	return exists
}
