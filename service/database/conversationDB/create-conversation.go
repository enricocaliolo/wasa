package conversationDB

import (
	"database/sql"
	"errors"
	"fmt"
	"wasa/service/shared/models"
)

func CreateConversation(db *sql.DB, members []int, name string) (models.Conversation, error) {
    if len(members) < 1 {
        return models.Conversation{}, errors.New("conversation must have at least one other member")
    }

    isGroup := len(members) > 1

    var conversation models.Conversation

    err := db.QueryRow(`
        INSERT INTO Conversation (name, is_group) 
        VALUES (?, ?) 
        RETURNING conversation_id, name, is_group, created_at
    `, name, isGroup).Scan(&conversation.ID, &conversation.Name, &conversation.Is_group, &conversation.Created_at)

    if err != nil {
        return conversation, fmt.Errorf("creating conversation: %w", err)
    }   

    for _, member := range members {
        _, err = db.Exec(`
            INSERT INTO ConversationParticipants (conversation_id, user_id)
            VALUES (?, ?)
        `, conversation.ID, member)
        if err != nil {
            return conversation, fmt.Errorf("adding participant %d: %w", member, err)
        }
    }

    return conversation, nil
}
