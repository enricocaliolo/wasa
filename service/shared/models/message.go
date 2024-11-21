package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID             int           `json:"id"`
	Content        []byte        `json:"content"`
	SentTime       time.Time     `json:"sent_time"`
	EditedTime     sql.NullTime  `json:"edited_time,omitempty"`
	DeletedTime    sql.NullTime  `json:"deleted_time,omitempty"`
	SenderID       sql.NullInt64 `json:"sender_id,omitempty"`
	ConversationID sql.NullInt64 `json:"conversation_id,omitempty"`
	RepliedTo      sql.NullInt64 `json:"replied_to,omitempty"`
	ForwardedFrom  sql.NullInt64 `json:"forwarded_from,omitempty"`
}
