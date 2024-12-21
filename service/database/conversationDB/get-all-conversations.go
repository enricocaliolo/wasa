package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)


func GetAllConversations(db *sql.DB, userID int) []models.Conversation {
    query := `SELECT 
        c.conversation_id,
        c.name,
        c.is_group,
        c.created_at
    FROM Conversation c
    JOIN ConversationParticipants cp ON c.conversation_id = cp.conversation_id
    WHERE cp.user_id = ?
    ORDER BY c.created_at DESC`

    rows, err := db.Query(query, userID)
    if err != nil {
        log.Printf("Error querying conversations: %v", err)
        return nil
    }
    defer rows.Close()

    var conversations []models.Conversation

    for rows.Next() {
        var conv models.Conversation
        var name sql.NullString
        
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