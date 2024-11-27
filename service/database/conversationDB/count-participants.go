package conversationDB

import "database/sql"

func CountParticipants(db *sql.DB, conversation_id int) (int, error) {
	query := `
		SELECT COUNT(*) FROM ConversationParticipants WHERE conversation_id = ?
	`

	var count int
	err := db.QueryRow(query, conversation_id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err

}
