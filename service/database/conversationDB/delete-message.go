package conversationDB

import (
	"database/sql"
)

func DeleteMessage(db *sql.DB, message_id int) (bool, error) {
	query := `UPDATE Message 
		SET deleted_time = CURRENT_TIMESTAMP 
		WHERE message_id = ?;`
	_, err := db.Exec(query, message_id)
	if err != nil {
		return false, err
	}
	return true, err
}
