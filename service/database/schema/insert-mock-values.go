package schema

import (
	"database/sql"
	"log"
)

func InsertMockValues(db *sql.DB) {
	sqlStatements := []string{
		`INSERT INTO User (username)
VALUES  ('enrico'),
        ('manu'),
        ('teste');`,

		`INSERT INTO "Conversation" (conversation_id, name, is_group)
VALUES  (1, "teste", true),
        (2, "teste", false);`,

		`INSERT INTO "ConversationParticipants" (conversation_id, user_id)
VALUES  (1, 1),
        (1, 2),
        (1, 3),
        (2, 1),
        (2, 3);`,

		`INSERT INTO "Message" (content, sender_id, conversation_id)
VALUES  ('oi, tudo bem?', 1, 1),
        ('tudo bem, e vc?', 2, 1),
        ('tudo bem tbm', 1, 1);
	`,
	}

	for _, sqlStatement := range sqlStatements {
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
