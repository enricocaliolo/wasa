package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func GetAllConversations(db *sql.DB, user_id int) []models.Conversation {
	query := `
        SELECT 
            c.conversation_id,
            c.name,
			c.is_group
        FROM Conversation c
        JOIN ConversationParticipants cp ON c.conversation_id = cp.conversation_id
        WHERE cp.user_id = ?
        ORDER BY c.created_at DESC`

	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conversation models.Conversation
		err := rows.Scan(
			&conversation.ID,
			&conversation.Name,
			&conversation.Is_group,
		)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		conversations = append(conversations, conversation)
	}

	if err = rows.Err(); err != nil {
		return nil
	}

	return conversations

}
