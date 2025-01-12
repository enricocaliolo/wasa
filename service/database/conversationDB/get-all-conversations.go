package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func GetAllConversations(db *sql.DB, userID int) ([]models.Conversation, error) {
	// First query to get conversations
	conversationsQuery := `
    SELECT 
        c.conversation_id,
        CASE 
            WHEN c.is_group = FALSE THEN (
                SELECT u.username 
                FROM User u 
                INNER JOIN ConversationParticipants cp2 ON u.user_id = cp2.user_id 
                WHERE cp2.conversation_id = c.conversation_id 
                AND cp2.user_id != ?
            )
            ELSE c.name 
        END as display_name,
		c.photo,
        c.is_group,
        c.created_at
    FROM Conversation c
    INNER JOIN ConversationParticipants cp ON c.conversation_id = cp.conversation_id
    WHERE cp.user_id = ?`

	rows, err := db.Query(conversationsQuery, userID, userID)
	if err != nil {
		log.Printf("Error querying conversations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation

	for rows.Next() {
		var conv models.Conversation
		var name string

		err := rows.Scan(
			&conv.ID,
			&name,
			&conv.Photo,
			&conv.Is_group,
			&conv.Created_at,
		)
		if err != nil {
			return nil, err
		}

		conv.Name = name

		participantsQuery := `
            SELECT 
                u.user_id,
                u.username,
                u.icon
            FROM ConversationParticipants cp
            JOIN User u ON cp.user_id = u.user_id
            WHERE cp.conversation_id = ?`

		participantRows, err := db.Query(participantsQuery, conv.ID)
		if err != nil {
			return nil, err
		}
		defer participantRows.Close()

		var participants []models.User
		for participantRows.Next() {
			var user models.User
			err := participantRows.Scan(
				&user.ID,
				&user.Username,
				&user.Icon,
			)
			if err != nil {
				return nil, err
			}
			participants = append(participants, user)
		}

		conv.Participants = participants
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}
