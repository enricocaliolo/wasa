package messagesdb

import (
	"database/sql"
	"wasa/service/shared/models"
)

func GetMessage(db *sql.DB, message_id, conversation_id int) (models.Message, error) {
	statement, _ := db.Prepare("SELECT message_id from Message WHERE message_id = ? AND conversation_id = ?")

	defer statement.Close()
	var message models.Message
	_ = statement.QueryRow(message_id, conversation_id).Scan(&message.Content, message.ContentType, message.ConversationID)

	return message, nil
}
