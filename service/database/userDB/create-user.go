package userDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func CreateUser(db *sql.DB, username string) (models.User, error) {
	user := models.User{}
	err := db.QueryRow("INSERT INTO User(username) VALUES (?) RETURNING user_id, username, icon, created_at",
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Icon,
		&user.Created_at,
	)
	if err != nil {
		return user, err
	}
	return user, nil

}
