package models

import "database/sql"

type User struct {
	ID         int
	Username   string
	Icon       sql.NullString
	Created_at sql.NullString
}
