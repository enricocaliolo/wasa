package conversationDB

import "database/sql"

func IsUserInConversation(db *sql.DB, userID int, conversationID int) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM ConversationParticipants
            WHERE conversation_id = ? AND user_id = ?
        )`

	var exists bool
	err := db.QueryRow(query, conversationID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
