package conversationDB

import (
	"database/sql"
	"fmt"
	"wasa/service/shared/models"
)

func SendMessage(db *sql.DB, message models.Message) (*models.Message, error) {
    query := `
    INSERT INTO Message (
        content,
        content_type,
        sender_id,
        conversation_id
    ) VALUES (?, ?, ?, ?)
    RETURNING 
        message_id,
        content,
        content_type,
        sent_time,
        edited_time,
        deleted_time,
        sender_id,
        conversation_id`
    
    var insertedMessage models.Message
    var senderID int
    
    err := db.QueryRow(query, 
        message.Content,
        message.ContentType,
        message.Sender.ID,
        message.ConversationID,
    ).Scan(
        &insertedMessage.ID,
        &insertedMessage.Content,
        &insertedMessage.ContentType,
        &insertedMessage.SentTime,
        &insertedMessage.EditedTime,
        &insertedMessage.DeletedTime,
        &senderID,
        &insertedMessage.ConversationID,
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
        return nil, fmt.Errorf("fetching sender info: %w", err)
    }
    
    return &insertedMessage, nil
}