package messagesdb

import (
	"database/sql"
	"strings"
)

func GetMessageSeenStatus(db *sql.DB, messageIDs []int) (map[int][]int, error) {
	if len(messageIDs) == 0 {
		return make(map[int][]int), nil
	}

	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "SELECT message_id, user_id FROM MessageSeen WHERE message_id IN (" +
		strings.Join(placeholders, ",") + ") ORDER BY message_id, seen_at"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seenStatus := make(map[int][]int)
	for _, messageID := range messageIDs {
		seenStatus[messageID] = []int{}
	}

	for rows.Next() {
		var messageID, userID int
		if err := rows.Scan(&messageID, &userID); err != nil {
			return nil, err
		}
		seenStatus[messageID] = append(seenStatus[messageID], userID)
	}

	return seenStatus, nil
}
