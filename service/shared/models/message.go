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
	IsForwarded  bool `json:"-"`
	Sender 	   User          `json:"user"`
	RepliedToMessage *Message    `json:"-"`
    Reactions      []Reaction    `json:"-"`
}

func (m Message) MarshalJSON() ([]byte, error) {
    return json.Marshal(&struct {
        ID              int        `json:"message_id"`
        Content         string     `json:"content"`
        ContentType     string     `json:"content_type"`
        SentTime        time.Time  `json:"sent_time"`
        EditedTime      *time.Time `json:"edited_time,omitempty"`
        DeletedTime     *time.Time `json:"deleted_time,omitempty"`
        SenderID        int        `json:"-"`
        ConversationID  int        `json:"conversation_id"`
        RepliedTo       *int64     `json:"-"`
        IsForwarded   bool     `json:"is_forwarded,omitempty"`
        Sender          User       `json:"sender"`
        RepliedToMessage *struct {  
            ID          int    `json:"message_id"`
            Content     string `json:"content"`
            ContentType string `json:"content_type"`
        } `json:"replied_to_message,omitempty"`
        Reactions       []Reaction `json:"reactions,omitempty"`
    }{
        ID:             m.ID,
        Content:        string(m.Content),
        ContentType:    m.ContentType,
        SentTime:       m.SentTime,
        EditedTime:     helper.NullTimeToPtr(m.EditedTime),
        DeletedTime:    helper.NullTimeToPtr(m.DeletedTime),
        RepliedTo:      helper.NullInt64ToPtr(m.RepliedTo),
        IsForwarded:  m.IsForwarded,
        SenderID:       m.senderID,
        ConversationID: m.ConversationID,
        Sender:         m.Sender,
        Reactions:     m.Reactions,
        RepliedToMessage: func() *struct {
            ID          int    `json:"message_id"`
            Content     string `json:"content"`
            ContentType string `json:"content_type"`
        } {
            if m.RepliedToMessage == nil {
                return nil
            }
            return &struct {
                ID          int    `json:"message_id"`
                Content     string `json:"content"`
                ContentType string `json:"content_type"`
            }{
                ID:          m.RepliedToMessage.ID,
                Content:     string(m.RepliedToMessage.Content),
                ContentType: m.RepliedToMessage.ContentType,
            }
        }(),
    })
}

