package models

import (
	"database/sql"
	"encoding/json"
)

type User struct {
	ID         int
	Username   string
	Icon       sql.NullString
	Created_at sql.NullString
}

func (c *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID         int    `json:"user_id"`
		Username   string `json:"username"`
		Icon       string `json:"icon"`
		Created_at string `json:"created_at"`
	}{
		ID:         c.ID,
		Username:   c.Username,
		Icon:       c.Icon.String,
		Created_at: c.Created_at.String,
	})
}
