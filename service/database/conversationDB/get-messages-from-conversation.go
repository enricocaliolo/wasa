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
            m.replied_to,
            m.forwarded_from,
            u.user_id,
            u.username,
            u.icon
        FROM Message m
        JOIN User u ON m.sender_id = u.user_id
        WHERE m.conversation_id = ?
        ORDER BY m.sent_time ASC`

    rows, err := db.Query(query, conversation_id)
    if err != nil {
        log.Printf("Error querying messages: %v", err)
        return nil
    }
    defer rows.Close()

    var messages []models.Message
    for rows.Next() {
        var msg models.Message
        var sender models.User
        var content []byte // since content is BLOB in the schema
        var editedTime, deletedTime sql.NullTime
        var repliedTo, forwardedFrom sql.NullInt64

        err := rows.Scan(
            &msg.ID,
            &content,
            &msg.ContentType,
            &msg.SentTime,
            &editedTime,
            &deletedTime,
            &repliedTo,
            &forwardedFrom,
            &sender.ID,
            &sender.Username,
            &sender.Icon,
        )
        if err != nil {
            log.Printf("Error scanning message: %v", err)
            continue
        }

        msg.Content = content
        msg.EditedTime = editedTime
        msg.DeletedTime = deletedTime
        msg.RepliedTo = repliedTo
        msg.ForwardedFrom = forwardedFrom
        msg.ConversationID = conversation_id
        msg.Sender = sender

        messages = append(messages, msg)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating messages: %v", err)
        return nil
    }

    return messages
}