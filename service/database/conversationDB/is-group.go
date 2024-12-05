package conversationDB

import (
	"database/sql"
)

func IsGroup(db *sql.DB, conversation_id int) (bool, error) {
	var isGroup bool
	err := db.QueryRow(`
        SELECT is_group FROM Conversation 
        WHERE conversation_id = ?
    `, conversation_id).Scan(&isGroup)
	return isGroup, err
}
