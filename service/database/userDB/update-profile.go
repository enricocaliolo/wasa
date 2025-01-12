package userDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func UpdateUsername(db *sql.DB, user models.User) bool {
	if isUsernameTaken(db, user.Username) {
		return false
	}

	statement, err := db.Prepare(`
    UPDATE User
    SET username = COALESCE(?, username)
    WHERE user_id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	_, err = statement.Exec(user.Username, user.ID)
	return err == nil
}

func UpdatePhoto(db *sql.DB, userId int, imageData []byte) bool {
	statement, err := db.Prepare(`
    UPDATE User
    SET icon = COALESCE(?, icon)
    WHERE user_id = ?
	`)
	if err != nil {
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(imageData, userId)
	return err == nil
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
