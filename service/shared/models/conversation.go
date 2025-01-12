package models

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Conversation struct {
	ID           int       `json:"conversation_id"`
	Name         string    `json:"name"`
	Photo        []byte    `json:"-"`
	Is_group     bool      `json:"is_group"`
	Created_at   time.Time `json:"created_at"`
	Messages     []Message `json:"messages"`
	Participants []User    `json:"participants"`
}

func (c *Conversation) MarshalJSON() ([]byte, error) {
	type Alias Conversation
	var photoStr string
	if len(c.Photo) > 0 {
		photoStr = base64.StdEncoding.EncodeToString(c.Photo)
	}

	return json.Marshal(&struct {
		*Alias
		Photo string `json:"photo,omitempty"`
	}{
		Alias: (*Alias)(c),
		Photo: photoStr,
	})
}
