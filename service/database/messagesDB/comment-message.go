package messagesdb

import (
	"database/sql"
)

func CommentMessage(db *sql.DB, user_id int, message_id int, reaction []byte) (bool, error) {
	query := `
        INSERT INTO Reactions (message_id, user_id, reaction)
        VALUES (?, ?, ?)
    `
	_, err := db.Exec(query, message_id, user_id, reaction)

	if err != nil {
		return false, err
	}

	return true, nil
}
