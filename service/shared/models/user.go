package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"time"
	"wasa/service/shared/helper"
)

type User struct {
	ID         int          `json:"user_id"`
	Username   string       `json:"username"`
	Icon       []byte       `json:"-"`
	Created_at sql.NullTime `json:"-"`
}

func (u User) MarshalJSON() ([]byte, error) {
	var iconStr string
	if len(u.Icon) > 0 {
		iconStr = base64.StdEncoding.EncodeToString(u.Icon)
	}

	return json.Marshal(struct {
		ID         int        `json:"user_id"`
		Username   string     `json:"username"`
		Icon       string     `json:"icon,omitempty"`
		Created_at *time.Time `json:"created_at,omitempty"`
	}{
		ID:         u.ID,
		Username:   u.Username,
		Icon:       iconStr,
		Created_at: helper.NullTimeToPtr(u.Created_at),
	})
}
