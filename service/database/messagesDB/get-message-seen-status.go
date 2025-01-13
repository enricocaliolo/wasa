package messagesdb

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetMessageSeenStatus(db *sql.DB, messageIDs []int) (map[int][]int, error) {
	if len(messageIDs) == 0 {
		return make(map[int][]int), nil
	}

	args := make([]interface{}, len(messageIDs))
	placeholders := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		args[i] = id
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
        SELECT message_id, user_id
        FROM MessageSeen
        WHERE message_id IN (%s)
        ORDER BY message_id, seen_at
    `, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying message seen status: %v", err)
	}
	defer rows.Close()

	seenStatus := make(map[int][]int)

	for _, messageID := range messageIDs {
		seenStatus[messageID] = []int{}
	}

	for rows.Next() {
		var messageID, userID int
		if err := rows.Scan(&messageID, &userID); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		seenStatus[messageID] = append(seenStatus[messageID], userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return seenStatus, nil
}
