package conversationDB

import (
	"database/sql"
)

func IsMessageFromUser(db *sql.DB, messageID int, userID int) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM Message 
            WHERE message_id = ? AND sender_id = ?
        )
    `

	var exists bool
	err := db.QueryRow(query, messageID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
