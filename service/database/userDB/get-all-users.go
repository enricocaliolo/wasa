package userDB

import (
	"database/sql"
	"wasa/service/shared/models"
)

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Icon, &user.Created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
