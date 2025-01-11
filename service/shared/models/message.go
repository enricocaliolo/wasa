package models

import (
	"database/sql"
	"encoding/base64"
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

    type messageAlias struct {
        ID              int        `json:"message_id"`
        Content         string     `json:"content"`
        ContentType     string     `json:"content_type"`
        SentTime        time.Time  `json:"sent_time"`
        EditedTime      *time.Time `json:"edited_time,omitempty"`
        DeletedTime     *time.Time `json:"deleted_time,omitempty"`
        SenderID        int        `json:"-"`
        ConversationID  int        `json:"conversation_id"`
        RepliedTo       *int64     `json:"-"`
        IsForwarded     bool       `json:"is_forwarded,omitempty"`
        Sender          User       `json:"sender"`
        RepliedToMessage *struct {  
            ID          int    `json:"message_id"`
            Content     string `json:"content"`
            ContentType string `json:"content_type"`
        } `json:"replied_to_message,omitempty"`
        Reactions       []Reaction `json:"reactions,omitempty"`
    }

    var content string
    if m.ContentType == "image" {
        content = base64.StdEncoding.EncodeToString(m.Content)
    } else {
        content = string(m.Content)
    }

    return json.Marshal(&messageAlias{
        ID:             m.ID,
        Content:        content,
        ContentType:    m.ContentType,
        SentTime:       m.SentTime,
        EditedTime:     helper.NullTimeToPtr(m.EditedTime),
        DeletedTime:    helper.NullTimeToPtr(m.DeletedTime),
        RepliedTo:      helper.NullInt64ToPtr(m.RepliedTo),
        IsForwarded:    m.IsForwarded,
        SenderID:       m.senderID,
        ConversationID: m.ConversationID,
        Sender:         m.Sender,
        Reactions:      m.Reactions,
        RepliedToMessage: func() *struct {
            ID          int    `json:"message_id"`
            Content     string `json:"content"`
            ContentType string `json:"content_type"`
        } {
            if m.RepliedToMessage == nil {
                return nil
            }
            
            var repliedContent string
            if m.RepliedToMessage.ContentType == "image" {
                repliedContent = base64.StdEncoding.EncodeToString(m.RepliedToMessage.Content)
            } else {
                repliedContent = string(m.RepliedToMessage.Content)
            }
            
            return &struct {
                ID          int    `json:"message_id"`
                Content     string `json:"content"`
                ContentType string `json:"content_type"`
            }{
                ID:          m.RepliedToMessage.ID,
                Content:     repliedContent,
                ContentType: m.RepliedToMessage.ContentType,
            }
        }(),
    })
}

