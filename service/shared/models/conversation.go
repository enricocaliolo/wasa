package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Conversation struct {
    ID                        int                       `json:"conversation_id"`
    Name                      string           `json:"name"`
    Is_group                  bool             `json:"is_group"`
    Created_at               time.Time           `json:"created_at"`
    Messages                 []Message                `json:"messages"`
    ConversationParticipant []User `json:"participants"`
}

func (c *Conversation) MarshalJSON() ([]byte, error) {
    return json.Marshal(&struct {
        ID                        int                       `json:"conversation_id"`
        Name                      string                    `json:"name"`
        IsGroup                   bool                      `json:"is_group"`
        CreatedAt                 time.Time                 `json:"created_at"`
        Messages                  []Message                 `json:"messages,omitempty"`
        ConversationParticipant   []User    `json:"participants,omitempty"`
    }{
        ID:                      c.ID,
        Name:                    c.Name,
        IsGroup:                 c.Is_group,          
        CreatedAt:               c.Created_at, 
        Messages:                c.Messages,
        ConversationParticipant: c.ConversationParticipant,
    })
}

func NullStringToPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}