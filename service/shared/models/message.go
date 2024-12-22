package models

import (
	"database/sql"
	"encoding/json"
	"time"
	"wasa/service/shared/helper"
)

type Message struct {
	ID             int           `json:"-"`
	Content        []byte        `json:"-"`
	ContentType    string        `json:"content_type"`
	SentTime       time.Time     `json:"sent_time"`
	EditedTime     sql.NullTime  `json:"-"`
	DeletedTime    sql.NullTime  `json:"-"`
	senderID       int           `json:"-"`
	ConversationID int           `json:"conversation_id"`
	RepliedTo      sql.NullInt64 `json:"-"`
	ForwardedFrom  sql.NullInt64 `json:"-"`
	Sender 	   User          `json:"user"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID             int        `json:"message_id"`
		Content        string     `json:"content"`
		ContentType    string     `json:"content_type"`
		SentTime       time.Time  `json:"sent_time"`
		EditedTime     *time.Time `json:"edited_time,omitempty"`
		DeletedTime    *time.Time `json:"deleted_time,omitempty"`
		SenderID       int        `json:"-"`
		ConversationID int        `json:"conversation_id"`
		RepliedTo      *int64     `json:"replied_to,omitempty"`
		ForwardedFrom  *int64     `json:"forwarded_from,omitempty"`
		Sender		 User       `json:"sender"`
	}{
		ID:             m.ID,
		Content:        string(m.Content),
		ContentType:    m.ContentType,
		SentTime:       m.SentTime,
		EditedTime:     helper.NullTimeToPtr(m.EditedTime),
		DeletedTime:    helper.NullTimeToPtr(m.DeletedTime),
		RepliedTo:      helper.NullInt64ToPtr(m.RepliedTo),
		ForwardedFrom:  helper.NullInt64ToPtr(m.ForwardedFrom),
		SenderID:       m.senderID,
		ConversationID: m.ConversationID,
		Sender: 	   	m.Sender,
	})
}

