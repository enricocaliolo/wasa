package conversationDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func GetConversation(db *sql.DB, conversation_id int) (*models.Conversation, error) {
	query := `
        SELECT conversation.conversation_id, conversation.name, conversation.photo, conversation.is_group, conversation.created_at 
        FROM Conversation conversation
        WHERE conversation.conversation_id = ?
    `
	conversation := &models.Conversation{}
	err := db.QueryRow(query, conversation_id).Scan(
		&conversation.ID,
		&conversation.Name,
		&conversation.Photo,
		&conversation.Is_group,
		&conversation.Created_at,
	)
	if err != nil {
		return nil, err
	}

	participantsQuery := `
        SELECT user.user_id, user.username, user.icon
        FROM User user
        JOIN ConversationParticipants cp ON cp.user_id = user.user_id
        WHERE cp.conversation_id = ?
    `
	participantsRows, err := db.Query(participantsQuery, conversation_id)
	if err != nil {
		return nil, err
	}
	defer participantsRows.Close()

	conversation.Participants = make([]models.User, 0)
	for participantsRows.Next() {
		var user models.User
		if err := participantsRows.Scan(&user.ID, &user.Username, &user.Icon); err != nil {
			return nil, err
		}
		conversation.Participants = append(conversation.Participants, user)
	}
	if err = participantsRows.Err(); err != nil {
		return nil, err
	}

	return conversation, nil
}
