package conversationDB

import (
	"database/sql"
	"fmt"
	"wasa/service/shared/models"
)

func ReplyToMessage(db *sql.DB, message models.Message) (*models.Message, error) {
    query := `
    INSERT INTO Message (
        content,
        content_type,
        sender_id,
        conversation_id,
        replied_to
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
        replied_to`
    
    var insertedMessage models.Message
    var senderID int
    
    err := db.QueryRow(query, 
        message.Content,
        message.ContentType,
        message.Sender.ID,
        message.ConversationID,
        message.RepliedTo.Int64,
    ).Scan(
        &insertedMessage.ID,
        &insertedMessage.Content,
        &insertedMessage.ContentType,
        &insertedMessage.SentTime,
        &insertedMessage.EditedTime,
        &insertedMessage.DeletedTime,
        &senderID,
        &insertedMessage.ConversationID,
        &insertedMessage.RepliedTo,
    )
    
    if err != nil {
        return nil, fmt.Errorf("inserting message: %w", err)
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

    if insertedMessage.RepliedTo.Valid {
        insertedMessage.RepliedToMessage = &models.Message{}
        err = db.QueryRow(`
            SELECT message_id, content, content_type
            FROM Message
            WHERE message_id = ?
        `, insertedMessage.RepliedTo.Int64).Scan(
            &insertedMessage.RepliedToMessage.ID,
            &insertedMessage.RepliedToMessage.Content,
            &insertedMessage.RepliedToMessage.ContentType,
        )
        if err != nil {
            return nil, fmt.Errorf("fetching replied-to message: %w", err)
        }
    }
    
    return &insertedMessage, nil
}