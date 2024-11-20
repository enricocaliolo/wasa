package schema

import (
	"database/sql"
	"log"
)

func InsertMockValues(db *sql.DB) {
	sqlStatements := []string{
		`INSERT INTO User (user_id, username)
VALUES  (1, 'enrico'),
        (2, 'manu'),
        (3, 'teste');`,

		`INSERT INTO "Conversation" (conversation_id)
VALUES  (1),
        (2);`,

		`INSERT INTO "ConversationParticipants" (conversation_id, user_id)
VALUES  (1, 1),
        (1, 2),
        (1, 3),
        (2, 1),
        (2, 3);`,

		`INSERT INTO "Message" (message_id, content, sender_id, conversation_id)
VALUES  (1, 'oi, tudo bem?', 1, 1),
        (2, 'tudo bem, e vc?', 2, 1),
        (3, 'tudo bem tbm', 1, 1);
	`,
	}

	for _, sqlStatement := range sqlStatements {
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
