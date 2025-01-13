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

	isGroup := len(members) > 2
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
	participantsQuery := `
        SELECT 
            u.user_id,
            u.username,
            u.icon
        FROM ConversationParticipants cp
        JOIN User u ON cp.user_id = u.user_id
        WHERE cp.conversation_id = ?`

	participantRows, err := db.Query(participantsQuery, conversation.ID)
	if err != nil {
		return conversation, fmt.Errorf("querying participants: %w", err)
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
			return conversation, fmt.Errorf("scanning participant: %w", err)
		}
		participants = append(participants, user)
	}

	if err = participantRows.Err(); err != nil {
		return conversation, fmt.Errorf("iterating participants: %w", err)
	}

	conversation.Participants = participants

	return conversation, nil
}
