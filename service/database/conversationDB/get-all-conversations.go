package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)


func GetAllConversations(db *sql.DB, userID int) []models.Conversation {
    query := `SELECT 
    c.conversation_id,
    CASE 
        WHEN c.name IS NULL AND c.is_group = FALSE THEN 
            (SELECT u.username 
             FROM ConversationParticipants cp 
             JOIN User u ON cp.user_id = u.user_id 
             WHERE cp.conversation_id = c.conversation_id 
               AND cp.user_id != 1 
             LIMIT 1)
        ELSE c.name 
    END AS conversation_name,
    c.is_group,
    c.created_at AS conversation_created_at,
    m.message_id,
    m.content,
    m.content_type,
    m.sent_time,
    m.sender_id,
    u.username AS sender_username,
    u.icon AS sender_icon
    FROM Conversation c
    JOIN ConversationParticipants cp ON c.conversation_id = cp.conversation_id
    LEFT JOIN Message m ON c.conversation_id = m.conversation_id
    LEFT JOIN User u ON m.sender_id = u.user_id
    WHERE cp.user_id = 1
    ORDER BY c.conversation_id, m.sent_time;`

    rows, err := db.Query(query, userID)
    if err != nil {
        log.Fatal(err)
        return nil
    }
    defer rows.Close()

    conversations := make(map[int]*models.Conversation)

for rows.Next() {
    var (
        conversationId int
        conversationName sql.NullString
        isGroup bool
        conversation_created_at sql.NullString
        messageId int
        content []byte
        contentType string
        sentTime sql.NullTime
        senderID int
        senderName sql.NullString
        senderIcon sql.NullString
    )

    err := rows.Scan(
        &conversationId,
        &conversationName,
        &isGroup,
        &conversation_created_at,
        &messageId,
        &content,
        &contentType,
        &sentTime,
        &senderID,
        &senderName,
        &senderIcon,
    )
    if err != nil {
        log.Fatal(err)
        return nil
    }

    conv, exists := conversations[conversationId]
    if !exists {
        conv = &models.Conversation{
            ID:         conversationId,
            Name:       conversationName,
            Is_group:   sql.NullBool{Bool: isGroup, Valid: true},
            Created_at: conversation_created_at,
            Messages:   make([]models.Message, 0),
        }
        conversations[conversationId] = conv
    }

    if sentTime.Valid {
        message := models.Message{
            ID:             messageId,
            Content:        content,
            ContentType:    "text",
            SentTime:       sentTime.Time,
            SenderID:       senderID,
            ConversationID: conversationId,
            Sender: models.User{
                ID:      senderID,
                Username: senderName.String,
                Icon:   senderIcon,
             },
        }
        conv.Messages = append(conv.Messages, message)
    }
}

result := make([]models.Conversation, 0, len(conversations))
for _, conv := range conversations {
    result = append(result, *conv)
}

return result

}