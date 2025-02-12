package conversationDB

import (
	"database/sql"
)

func DeleteConversation(db *sql.DB, conversation_id int) (bool, error) {
	query := `UPDATE Message 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE conversation_id = ?;`
	_, err := db.Exec(query, conversation_id)
	if err != nil {
		return false, err
	}
	return true, err
}
