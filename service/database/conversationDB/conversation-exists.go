package conversationDB

import (
	"database/sql"
)

func ConversationExists(db *sql.DB, conversationID int) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM Conversation 
            WHERE conversation_id = ?
        )
    `

	err := db.QueryRow(query, conversationID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
