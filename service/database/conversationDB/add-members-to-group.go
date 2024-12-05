package conversationDB

import (
	"database/sql"
	"fmt"
)

func AddGroupMembers(db *sql.DB, conversation_id int, members []int) error {
	for _, member := range members {
		_, err := db.Exec(`
            INSERT INTO ConversationParticipants (conversation_id, user_id)
            VALUES (?, ?)
        `, conversation_id, member)
		if err != nil {
			return fmt.Errorf("adding member %d: %w", member, err)
		}
	}
	return nil
}
