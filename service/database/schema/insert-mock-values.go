package schema

import (
	"database/sql"
	"log"
)

func InsertMockValues(db *sql.DB) {
	sqlStatements := []string{
		// `INSERT INTO "User" (username) VALUES
		// ('enrico'),
		// ('jane_smith'),
		// ('bob_wilson'),
		// ('alice_green'),
		// ('charlie_brown'),
		// ('diana_prince'),
		// ('bruce_wayne'),
		// ('peter_parker'),
		// ('tony_stark'),
		// ('steve_rogers');`,

		// `INSERT INTO "Conversation" (name, is_group) VALUES
		// (NULL, FALSE),
		// (NULL, FALSE),
		// ('Project Team', TRUE),
		// ('Family Group', TRUE),
		// (NULL, FALSE),
		// ('Gaming Squad', TRUE),
		// ('Book Club', TRUE),
		// (NULL, FALSE),
		// ('Tech Team', TRUE),
		// (NULL, FALSE);`,

		// `INSERT INTO "ConversationParticipants" (conversation_id, user_id) VALUES
		// (1, 1), (1, 2),
		// (2, 3), (2, 4),
		// (3, 1), (3, 2), (3, 3), (3, 4),
		// (4, 5), (4, 6), (4, 7), (4, 8),
		// (5, 5), (5, 6),
		// (6, 7), (6, 8), (6, 9), (6, 10),
		// (7, 1), (7, 3), (7, 5), (7, 7),
		// (8, 8), (8, 9),
		// (9, 8), (9, 9), (9, 10),
		// (10, 7), (10, 10);`,

		// `INSERT INTO "Message" (content, content_type, sender_id, conversation_id, replied_to, is_forwarded) VALUES
		// ('Hey, how are you?', 'text', 1, 1, NULL, FALSE),
		// ('Im doing great, thanks for asking!', 'text', 2, 1, 1, FALSE),
		// ('Project meeting at 3 PM tomorrow', 'text', 1, 3, NULL, FALSE),
		// ('Ill be there', 'text', 2, 3, 3, FALSE),
		// ('Can anyone help with the database issue?', 'text', 3, 3, NULL, FALSE),
		// ('Sure, whats the problem?', 'text', 4, 3, 5, FALSE),
		// ('Dinner at 7?', 'text', 5, 4, NULL, FALSE),
		// ('Perfect timing!', 'text', 6, 4, 7, FALSE),
		// ('New game release this weekend', 'text', 7, 6, NULL, FALSE),
		// ('Count me in for the launch', 'text', 8, 6, 9, FALSE),
		// ('Has everyone read chapter 5?', 'text', 1, 7, NULL, FALSE),
		// ('Yes, great plot twist!', 'text', 3, 7, 11, FALSE),
		// ('Project deadline extended to next week', 'text', 8, 9, NULL, FALSE),
		// ('Thanks for the update', 'text', 9, 9, 13, FALSE),
		// ('Meeting notes from today', 'text', 10, 9, NULL, FALSE),
		// ('New feature deployment tonight', 'text', 7, 10, NULL, FALSE),
		// ('All systems ready', 'text', 10, 10, 16, FALSE),
		// ('Remember to update the documentation', 'text', 8, 9, NULL, FALSE),
		// ('Already on it', 'text', 9, 9, 18, FALSE),
		// ('Good morning team!', 'text', 1, 3, NULL, FALSE);
		// `,

		// `
		// INSERT INTO "Reactions" (message_id, user_id, reaction) VALUES
		// (1, 2, '👍'),
		// (1, 3, '❤️'),
		// (2, 1, '😊'),
		// (3, 2, '👀'),
		// (3, 4, '✅'),
		// (7, 6, '🎉'),
		// (7, 8, '👍'),
		// (9, 8, '🎮'),
		// (9, 10, '🎮'),
		// (13, 9, '📝'),
		// (16, 10, '🚀');
		// `,
	}

	for _, sqlStatement := range sqlStatements {
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Printf("Error inserting mock values: %v", err)
		}
	}
}
