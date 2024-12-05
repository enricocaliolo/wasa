package conversationDB

import (
	"database/sql"
)

func UpdateGroupPhoto(db *sql.DB, conversation_id int, photo string) (bool, error) {
	_, err := db.Exec(`
        UPDATE Conversation 
        SET icon = ? 
        WHERE conversation_id = ? AND is_group = true
    `, photo, conversation_id)
	if err != nil {
		return false, err
	}
	return true, nil
}
