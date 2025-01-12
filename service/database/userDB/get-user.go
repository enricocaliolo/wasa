package userDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func GetUser(db *sql.DB, username string) (models.User, error) {
	statement, err := db.Prepare("SELECT * from User WHERE username = ?")
	user := models.User{}
	if err != nil {
		return user, err
	}
	defer statement.Close()
	err = statement.QueryRow(username).Scan(
		&user.ID, &user.Username, &user.Icon, &user.Created_at,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}
