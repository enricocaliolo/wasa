package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func GetConversationParticipants(db *sql.DB, conversationID int) []models.User {
	query := `
        SELECT 
            u.user_id,
            u.username,
            u.icon
        FROM ConversationParticipants cp
        JOIN User u ON cp.user_id = u.user_id
        WHERE cp.conversation_id = ?`

	rows, err := db.Query(query, conversationID)
	if err != nil {
		log.Printf("Error querying conversation participants: %v", err)
		return nil
	}
	defer rows.Close()

	var participants []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Icon,
		)
		if err != nil {
			log.Printf("Error scanning participant: %v", err)
			continue
		}

		participants = append(participants, user)
	}

	return participants
}
