package conversationDB

import (
	"database/sql"
)

func UpdateGroupName(db *sql.DB, conversation_id int, name string) (bool, error) {
	_, err := db.Exec(`
        UPDATE Conversation 
        SET name = ? 
        WHERE conversation_id = ? AND is_group = true
    `, name, conversation_id)
	if err != nil {
		return false, err
	}
	return true, nil
}
