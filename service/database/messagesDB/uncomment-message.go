package messagesdb

import (
	"database/sql"
)

func UncommentMessage(db *sql.DB, reaction_id int) (bool, error) {
	query := `DELETE FROM Reactions WHERE reaction_id = ?`
	_, err := db.Exec(query, reaction_id)
	if err != nil {
		return false, nil
	}
	return true, nil
}
