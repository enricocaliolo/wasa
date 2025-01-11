package models

import (
	"encoding/json"
)

type Reaction struct {
    ID        int    `json:"-"`
    MessageID int    `json:"-"`
    UserID    int    `json:"-"`
    Reaction  string `json:"-"`
    User      User   `json:"user"`
}

func (r Reaction) MarshalJSON() ([]byte, error) {
    return json.Marshal(&struct {
        ID        int    `json:"reaction_id"`
        MessageID int    `json:"message_id"`
        UserID    int    `json:"-"`
        Reaction  string `json:"reaction"`
        User      User   `json:"user"`
    }{
        ID:        r.ID,
        MessageID: r.MessageID,
        UserID:    r.UserID,
        Reaction:  r.Reaction,
        User:      r.User,
    })
}