package conversationDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func ForwardMessage(db *sql.DB, message models.Message) (*models.Message, error) {
	query := `
    INSERT INTO Message (
        content,
        content_type,
        sender_id,
        conversation_id,
        is_forwarded
    ) VALUES (?, ?, ?, ?, ?)
    RETURNING 
        message_id,
        content,
        content_type,
        sent_time,
        edited_time,
        deleted_time,
        sender_id,
        conversation_id,
        is_forwarded`

	var insertedMessage models.Message
	var senderID int

	err := db.QueryRow(query,
		message.Content,
		message.ContentType,
		message.Sender.ID,
		message.ConversationID,
		message.IsForwarded,
	).Scan(
		&insertedMessage.ID,
		&insertedMessage.Content,
		&insertedMessage.ContentType,
		&insertedMessage.SentTime,
		&insertedMessage.EditedTime,
		&insertedMessage.DeletedTime,
		&senderID,
		&insertedMessage.ConversationID,
		&insertedMessage.IsForwarded,
	)

	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`
        SELECT user_id, username, icon
        FROM User
        WHERE user_id = ?
    `, senderID).Scan(
		&insertedMessage.Sender.ID,
		&insertedMessage.Sender.Username,
		&insertedMessage.Sender.Icon,
	)

	if err != nil {
		return nil, err
	}

	return &insertedMessage, nil
}
