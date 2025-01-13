package messagesdb

import (
	"database/sql"
	"fmt"
)

func MarkMessagesSeen(db *sql.DB, userID int, messageIDs []int) error {
	fmt.Printf("Marking messages %v as seen by user %d\n", messageIDs, userID)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT OR REPLACE INTO MessageSeen (message_id, user_id, seen_at)
        VALUES (?, ?, CURRENT_TIMESTAMP)
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, messageID := range messageIDs {
		_, err = stmt.Exec(messageID, userID)
		if err != nil {
			return fmt.Errorf("error marking message %d as seen: %v", messageID, err)
		}
		fmt.Printf("Marked message %d as seen by user %d\n", messageID, userID)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	fmt.Printf("Successfully marked all messages as seen\n")
	return nil
}
