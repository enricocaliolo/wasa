package messagesdb

import (
	"database/sql"
	"wasa/service/shared/models"
)

func UncommentMessage(db *sql.DB, reaction_id int) (models.Reaction, error) {
	var reaction models.Reaction
	var user models.User

	query := `
        DELETE FROM Reactions 
        WHERE reaction_id = ? 
        RETURNING reaction_id, message_id, reaction, user_id
    `

	err := db.QueryRow(query, reaction_id).Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.Reaction,
		&user.ID,
	)

	if err != nil {
		return reaction, err
	}

	userQuery := `SELECT username, icon FROM "User" WHERE user_id = ?`
	err = db.QueryRow(userQuery, user.ID).Scan(&user.Username, &user.Icon)
	if err != nil {
		return reaction, err
	}

	reaction.User = user
	return reaction, nil
}
