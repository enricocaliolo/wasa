package models

import (
	"database/sql"
)

type Conversation struct {
	ID         int            `json:"id"`
	Name       sql.NullString `json:"Name"`
	Is_group   sql.NullBool   `json:"Is_group"`
	Created_at sql.NullString `json:"Created_at"`
	Messages   []Message      `json:"Messages"`
}
