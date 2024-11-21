package userDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func UpdateProfile(db *sql.DB, user models.User) string {
	statement, err := db.Prepare(`
    UPDATE User
    SET username = COALESCE(?, username),
        icon = COALESCE(?, icon)
    WHERE user_id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	_, err = statement.Exec(user.Username, user.Icon, user.ID)

	return user.Username
}
