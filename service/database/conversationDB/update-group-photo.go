package conversationDB

import (
	"database/sql"
)

func UpdateGroupPhoto(db *sql.DB, conversation_id int, photo []byte) (bool, error) {
	// result, err := db.Exec(`
	//     UPDATE Conversation
	//     SET photo = ?
	//     WHERE conversation_id = ? AND is_group = true
	// `, photo, conversation_id)

	// if err != nil {
	// 	return false, err
	// }

	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return false, err
	// }

	// return rowsAffected > 0, nil

	statement, err := db.Prepare(`
    UPDATE Conversation
    SET photo = COALESCE(?, photo)
    WHERE conversation_id = ?
	`)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	_, err = statement.Exec(photo, conversation_id)
	if err != nil {
		return false, err
	}
	return true, nil
}
