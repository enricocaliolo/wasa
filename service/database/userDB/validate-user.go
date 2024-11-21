package userDB

import (
	"database/sql"
)

func ValidateUser(db *sql.DB, id int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE user_id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return false
	}

	return true
}
