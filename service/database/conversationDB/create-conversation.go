package conversationDB

import (
	"database/sql"
	"errors"
	"fmt"
)

func CreateConversation(db *sql.DB, creator_id int, members []int) (int, error) {
	if len(members) < 1 {
		return 0, errors.New("conversation must have at least one other member")
	}

	isGroup := len(members) > 1

	result, err := db.Exec(`
        INSERT INTO Conversation (is_group) VALUES (?)
    `, isGroup)
	if err != nil {
		return 0, err
	}

	conversation_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Add creator and members
	participants := append([]int{creator_id}, members...)
	for _, member := range participants {
		_, err = db.Exec(`
            INSERT INTO ConversationParticipants (conversation_id, user_id)
            VALUES (?, ?)
        `, conversation_id, member)
		if err != nil {
			return 0, fmt.Errorf("adding participant %d: %w", member, err)
		}
	}

	return int(conversation_id), nil
}
