package userDB

import (
	"database/sql"
	"errors"
	"wasa/service/shared/models"
)

func UpdateUsername(db *sql.DB, user models.User) (bool, error) {
	if isUsernameTaken(db, user.Username) {
		return false, errors.New("Username is already taken")
	}

	statement, err := db.Prepare(`
    UPDATE User
    SET username = COALESCE(?, username)
    WHERE user_id = ?
	`)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	_, err = statement.Exec(user.Username, user.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdatePhoto(db *sql.DB, userId int, imageData []byte) (bool, error) {
	statement, err := db.Prepare(`
    UPDATE User
    SET icon = COALESCE(?, icon)
    WHERE user_id = ?
	`)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	_, err = statement.Exec(imageData, userId)
	if err != nil {
		return false, nil
	}
	return true, nil
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
