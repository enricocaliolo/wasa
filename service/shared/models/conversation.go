package models

import "database/sql"

type Conversation struct {
	ID         int
	Name       sql.NullString
	Is_group   sql.NullBool
	Created_at sql.NullString
}
