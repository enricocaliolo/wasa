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
            m.is_forwarded,
            u.user_id,
            u.username,
            u.icon,
            rm.message_id,
            rm.content,
            rm.content_type
        FROM Message m
        JOIN User u ON m.sender_id = u.user_id
        LEFT JOIN Message rm ON m.replied_to = rm.message_id
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
        var content []byte
        var editedTime, deletedTime sql.NullTime
        var repliedTo sql.NullInt64
        var isForwarded bool
        var repliedToMsg struct {
            ID          sql.NullInt64
            Content     []byte
            ContentType sql.NullString
        }

        err := rows.Scan(
            &msg.ID,
            &content,
            &msg.ContentType,
            &msg.SentTime,
            &editedTime,
            &deletedTime,
            &repliedTo,
            &isForwarded,
            &sender.ID,
            &sender.Username,
            &sender.Icon,
            &repliedToMsg.ID,
            &repliedToMsg.Content,
            &repliedToMsg.ContentType,
        )
        if err != nil {
            log.Printf("Error scanning message: %v", err)
            continue
        }

        msg.Content = content
        msg.EditedTime = editedTime
        msg.DeletedTime = deletedTime
        msg.RepliedTo = repliedTo
        msg.IsForwarded = isForwarded
        msg.ConversationID = conversation_id
        msg.Sender = sender

        // If there's a replied-to message, include its details
        if repliedToMsg.ID.Valid {
            msg.RepliedToMessage = &models.Message{
                ID:          int(repliedToMsg.ID.Int64),
                Content:     repliedToMsg.Content,
                ContentType: repliedToMsg.ContentType.String,
            }
        }

        messages = append(messages, msg)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating messages: %v", err)
        return nil
    }

    return messages
}