package models

import (
	"database/sql"
	"encoding/json"
	"time"
	"wasa/service/shared/helper"
)

type User struct {
	ID         int `json:"user_id"`
	Username   string `json:"username"`
	Icon       sql.NullString `json:"-"`
	Created_at sql.NullTime `json:"-"`
}

func (c *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID         int    `json:"user_id"`
		Username   string `json:"username"`
		Icon       *string `json:"icon,omitempty"`
		Created_at *time.Time `json:"created_at,omitempty"`
	}{
		ID:         c.ID,
		Username:   c.Username,
		Icon:       helper.NullStringToPtr(c.Icon),
		Created_at: helper.NullTimeToPtr(c.Created_at),
	})
}
