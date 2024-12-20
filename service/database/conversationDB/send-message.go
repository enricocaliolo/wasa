package conversationDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func SendMessage(db *sql.DB, message models.Message) (*models.Message, error) {
    query := `
    INSERT INTO Message (
        content,
        content_type,
        sender_id,
        conversation_id,
        replied_to,
        forwarded_from
    ) VALUES (?, ?, ?, ?, ?, ?)
    RETURNING message_id, content, content_type, sent_time, edited_time, 
              deleted_time, sender_id, conversation_id, replied_to, forwarded_from`
    
    var insertedMessage models.Message
    err := db.QueryRow(query, 
        message.Content,
        message.ContentType,
        message.SenderID,
        message.ConversationID,
        message.RepliedTo.Int64,
        message.ForwardedFrom.Int64,
    ).Scan(
        &insertedMessage.ID,
        &insertedMessage.Content,
        &insertedMessage.ContentType,
        &insertedMessage.SentTime,
        &insertedMessage.EditedTime,
        &insertedMessage.DeletedTime,
        &insertedMessage.SenderID,
        &insertedMessage.ConversationID,
        &insertedMessage.RepliedTo,
        &insertedMessage.ForwardedFrom,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &insertedMessage, nil
}
