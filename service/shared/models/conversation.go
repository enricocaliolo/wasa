package models

import (
	"database/sql"
	"encoding/json"
)

type Conversation struct {
    ID                        int                       `json:"id"`
    Name                      sql.NullString           `json:"-"`
    Is_group                  sql.NullBool             `json:"-"`
    Created_at               sql.NullString           `json:"-"`
    Messages                 []Message                `json:"messages"`
    ConversationParticipant []ConversationParticipant `json:"conversation_participants"`
}

func (c *Conversation) MarshalJSON() ([]byte, error) {
    return json.Marshal(&struct {
        ID                        int                       `json:"id"`
        Name                      string                    `json:"name"`
        Is_group                  bool                      `json:"is_group"`
        Created_at               string                    `json:"created_at"`
        Messages                 []Message                 `json:"messages"`
        ConversationParticipant []ConversationParticipant `json:"conversation_participants"`
    }{
        ID:                      c.ID,
        Name:                    c.Name.String,
        Is_group:                c.Is_group.Bool,
        Created_at:             c.Created_at.String,
        Messages:               c.Messages,
        ConversationParticipant: c.ConversationParticipant,
    })
}

type ConversationParticipant struct {
	User_id int `json:"user_id"`
	Joined_at string `json:"joined_at"`
    Name string `json:"name"`
}