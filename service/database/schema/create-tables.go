package schema

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {

	sqlStatements := []string{
		`DROP TABLE IF EXISTS "User";
		CREATE TABLE "User" (
		user_id INTEGER PRIMARY KEY,
		username VARCHAR(64) NOT NULL,
		icon VARCHAR(64),
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		`DROP TABLE IF EXISTS "Conversation";
		CREATE TABLE IF NOT EXISTS "Conversation" (
			conversation_id INTEGER PRIMARY KEY,
			"name" VARCHAR(64) DEFAULT '',
			is_group BOOLEAN,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		`DROP TABLE IF EXISTS "ConversationParticipants";
		CREATE TABLE IF NOT EXISTS "ConversationParticipants" (
			conversation_id INTEGER,
			user_id INTEGER,
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (conversation_id, user_id),
			FOREIGN KEY (conversation_id) REFERENCES "Conversation"(conversation_id),
			FOREIGN KEY (user_id) REFERENCES "User"(user_id)
		);`,

		`DROP TABLE IF EXISTS "Message";
		CREATE TABLE IF NOT EXISTS "Message" (
			message_id INTEGER PRIMARY KEY,
			content BLOB NOT NULL,
			content_type TEXT NOT NULL DEFAULT 'text',
			sent_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			edited_time TIMESTAMP,
			deleted_time TIMESTAMP,
			sender_id INTEGER,
			conversation_id INTEGER,
			replied_to INTEGER,
			forwarded_from INTEGER,
			FOREIGN KEY (sender_id) REFERENCES "User"(user_id),
			FOREIGN KEY (conversation_id) REFERENCES "Conversation"(conversation_id),
			FOREIGN KEY (replied_to) REFERENCES "Message"(message_id),
			FOREIGN KEY (forwarded_from) REFERENCES "Message"(message_id)
		);`,

		`DROP TABLE IF EXISTS "Reactions";
		CREATE TABLE IF NOT EXISTS "Reactions" (
			reaction_id INTEGER AUTO_INCREMENT PRIMARY KEY,
			message_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			reaction BLOB NOT NULL,
			FOREIGN KEY (message_id) REFERENCES "Message"(message_id),
			FOREIGN KEY (user_id) REFERENCES "User"(user_id)
		)`,
	}

	for _, sqlStatement := range sqlStatements {
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
