package conversationDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func DeleteMessage(db *sql.DB, message_id int) (models.Message, error) {
	query := `
        SELECT 
            m.message_id,
            m.content,
            m.content_type,
            m.sent_time,
            m.edited_time,
            m.deleted_time,
            m.replied_to,
            m.is_forwarded,
            m.conversation_id,
            u.user_id,
            u.username,
            u.icon
        FROM Message m
        JOIN User u ON m.sender_id = u.user_id
        WHERE m.message_id = ?;`

	var msg models.Message
	var sender models.User
	var content []byte
	var editedTime, deletedTime sql.NullTime
	var repliedTo sql.NullInt64

	_, err := db.Exec(`UPDATE Message SET deleted_time = CURRENT_TIMESTAMP WHERE message_id = ?`, message_id)
	if err != nil {
		return msg, err
	}
	err = db.QueryRow(query, message_id).Scan(
		&msg.ID,
		&content,
		&msg.ContentType,
		&msg.SentTime,
		&editedTime,
		&deletedTime,
		&repliedTo,
		&msg.IsForwarded,
		&msg.ConversationID,
		&sender.ID,
		&sender.Username,
		&sender.Icon,
	)

	if err != nil {
		return msg, err
	}

	msg.Content = content
	msg.EditedTime = editedTime
	msg.DeletedTime = deletedTime
	msg.RepliedTo = repliedTo
	msg.Sender = sender

	return msg, nil
}
