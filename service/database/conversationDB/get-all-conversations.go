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
        WHEN c.is_group THEN COALESCE(c.name, '')
        ELSE (
            SELECT u.username 
            FROM User u 
            JOIN ConversationParticipants cp2 ON u.user_id = cp2.user_id 
            WHERE cp2.conversation_id = c.conversation_id 
            AND cp2.user_id != ?
        )
    END as name,
    c.is_group,
    c.created_at
FROM Conversation c
JOIN ConversationParticipants cp ON c.conversation_id = cp.conversation_id
WHERE cp.user_id = ?
ORDER BY c.created_at DESC;`

    rows, err := db.Query(query, userID, userID)
    if err != nil {
        log.Printf("Error querying conversations: %v", err)
        return nil
    }
    defer rows.Close()

    var conversations []models.Conversation

    for rows.Next() {
        var conv models.Conversation
        var name string
        
        err := rows.Scan(
            &conv.ID,
            &name,
            &conv.Is_group,
            &conv.Created_at,
        )
        if err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }

        conv.Name = name
        conversations = append(conversations, conv)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil
    }

    return conversations
}