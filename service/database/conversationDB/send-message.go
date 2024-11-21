package conversationDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func SendMessage(db *sql.DB, message models.Message) (int, error) {

	query := `
        INSERT INTO Message (
            content,
			content_type,
            sender_id,
            conversation_id,
            replied_to,
            forwarded_from
        ) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, message.Content, message.ContentType, message.SenderID, message.ConversationID, message.RepliedTo, message.ForwardedFrom)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
