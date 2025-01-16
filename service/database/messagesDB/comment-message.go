package messagesdb

import (
	"database/sql"
	"wasa/service/shared/models"
)

func CommentMessage(db *sql.DB, user_id int, message_id int, reaction string) (models.Reaction, error) {
	var newReaction models.Reaction
	var user models.User

	err := db.QueryRow(`
        INSERT INTO Reactions (message_id, user_id, reaction)
        VALUES (?, ?, ?)
        RETURNING reaction_id, message_id, user_id, reaction
    `, message_id, user_id, reaction).Scan(
		&newReaction.ID,
		&newReaction.MessageID,
		&newReaction.UserID,
		&newReaction.Reaction,
	)

	if err != nil {
		return newReaction, err
	}

	err = db.QueryRow(`
        SELECT user_id, username, icon
        FROM User
        WHERE user_id = ?
    `, user_id).Scan(
		&user.ID,
		&user.Username,
		&user.Icon,
	)

	if err != nil {
		return newReaction, err
	}

	newReaction.User = user
	return newReaction, nil
}
