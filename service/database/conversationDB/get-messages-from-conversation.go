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
        
                if repliedToMsg.ID.Valid {
                    msg.RepliedToMessage = &models.Message{
                        ID:          int(repliedToMsg.ID.Int64),
                        Content:     repliedToMsg.Content,
                        ContentType: repliedToMsg.ContentType.String,
                    }
                }
        
                messages = append(messages, msg)
            }

    reactionQuery := `
        SELECT 
            r.reaction_id,
            r.message_id,
            r.user_id,
            r.reaction,
            u.user_id,
            u.username,
            u.icon
        FROM Reactions r
        JOIN User u ON r.user_id = u.user_id
        WHERE r.message_id IN (
            SELECT message_id 
            FROM Message 
            WHERE conversation_id = ?
        )`
        
    reactionRows, err := db.Query(reactionQuery, conversation_id)
    if err != nil {
        log.Printf("Error querying reactions: %v", err)
        return messages
    }
    defer reactionRows.Close()

    reactionMap := make(map[int][]models.Reaction)
    for reactionRows.Next() {
        var r models.Reaction
        var u models.User
        err := reactionRows.Scan(
            &r.ID,
            &r.MessageID,
            &r.UserID,
            &r.Reaction,
            &u.ID,
            &u.Username,
            &u.Icon,
        )
        if err != nil {
            log.Printf("Error scanning reaction: %v", err)
            continue
        }
        r.User = u
        reactionMap[r.MessageID] = append(reactionMap[r.MessageID], r)
    }

    for i := range messages {
        if reactions, ok := reactionMap[messages[i].ID]; ok {
            messages[i].Reactions = reactions
        }
    }

    return messages
}