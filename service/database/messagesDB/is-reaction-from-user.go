package messagesdb

import (
	"database/sql"
)

func IsReactionFromUser(db *sql.DB, reaction_id int, user_id int) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM Reactions 
            WHERE reaction_id = ? AND user_id = ?
        )
    `

	var exists bool
	err := db.QueryRow(query, reaction_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
