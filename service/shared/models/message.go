package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Message struct {
	ID             int           `json:"-"`
	Content        []byte        `json:"-"`
	ContentType    string        `json:"content_type"`
	SentTime       time.Time     `json:"sent_time"`
	EditedTime     sql.NullTime  `json:"-"`
	DeletedTime    sql.NullTime  `json:"-"`
	SenderID       int           `json:"sender_id"`
	ConversationID int           `json:"conversation_id"`
	RepliedTo      sql.NullInt64 `json:"-"`
	ForwardedFrom  sql.NullInt64 `json:"-"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID             int        `json:"message_id"`
		Content        string     `json:"content"`
		ContentType    string     `json:"content_type"`
		SentTime       time.Time  `json:"sent_time"`
		EditedTime     *time.Time `json:"edited_time,omitempty"`
		DeletedTime    *time.Time `json:"deleted_time,omitempty"`
		SenderID       int        `json:"sender_id"`
		ConversationID int        `json:"conversation_id"`
		RepliedTo      *int64     `json:"replied_to,omitempty"`
		ForwardedFrom  *int64     `json:"forwarded_from,omitempty"`
	}{
		ID:             m.ID,
		Content:        string(m.Content),
		ContentType:    m.ContentType,
		SentTime:       m.SentTime,
		EditedTime:     nullTimeToPtr(m.EditedTime),
		DeletedTime:    nullTimeToPtr(m.DeletedTime),
		SenderID:       m.SenderID,
		ConversationID: m.ConversationID,
		RepliedTo:      nullInt64ToPtr(m.RepliedTo),
		ForwardedFrom:  nullInt64ToPtr(m.ForwardedFrom),
	})
}

func nullInt64ToPtr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

func nullTimeToPtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
