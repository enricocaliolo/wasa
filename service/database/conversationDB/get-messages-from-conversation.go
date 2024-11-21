package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func GetMessagesFromConversation(db *sql.DB, conversation_id int) []models.Message {
	query := `
        SELECT 
            m.message_id,
            m.content,
			m.content_type,
            m.sent_time,
            m.edited_time,
            m.deleted_time,
            m.sender_id,
            m.replied_to,
            m.forwarded_from
        FROM Message m
        LEFT JOIN User u ON m.sender_id = u.user_id
        WHERE m.conversation_id = ?
        ORDER BY m.sent_time ASC`

	rows, err := db.Query(query, conversation_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.ContentType,
			&msg.SentTime,
			&msg.EditedTime,
			&msg.DeletedTime,
			&msg.SenderID,
			&msg.RepliedTo,
			&msg.ForwardedFrom,
		)
		if err != nil {
			log.Fatal(err)
		}
		msg.ConversationID = conversation_id
		messages = append(messages, msg)
	}

	return messages

}
