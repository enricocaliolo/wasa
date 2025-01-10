package conversationDB

import (
	"database/sql"
	"fmt"
)

func ConversationExists(db *sql.DB, conversationID int) (bool, error) {
    var exists bool
    query := `
        SELECT EXISTS(
            SELECT 1 
            FROM Conversation 
            WHERE conversation_id = ?
        )
    `
    
    err := db.QueryRow(query, conversationID).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("checking conversation existence: %w", err)
    }
    
    return exists, nil
}