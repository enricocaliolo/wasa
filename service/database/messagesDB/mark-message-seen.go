package messagesdb

import (
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

func MarkMessagesSeen(db *sql.DB, userID int, messageIDs []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logrus.Error("Error rolling back transaction: ", err)
		}
	}()

	stmt, err := tx.Prepare(`
        INSERT OR REPLACE INTO MessageSeen (message_id, user_id)
        VALUES (?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, messageID := range messageIDs {
		_, err = stmt.Exec(messageID, userID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
