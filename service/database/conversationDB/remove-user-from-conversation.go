package conversationDB

import (
	"database/sql"
)

func RemoveUserFromConversation(db *sql.DB, conversation_id int, user_id int) (bool, error) {
	query := `DELETE FROM ConversationParticipants WHERE conversation_id = ? AND user_id = ?`
	_, err := db.Exec(query, conversation_id, user_id)
	if err != nil {
		return false, err
	}
	return true, err
}
