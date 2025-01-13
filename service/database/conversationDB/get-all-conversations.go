package conversationDB

import (
	"database/sql"
	"log"
	"wasa/service/shared/models"
)

func GetAllConversations(db *sql.DB, userID int) ([]models.Conversation, error) {
	// First query to get conversations
	conversationsQuery := `
    SELECT 
        c.conversation_id,
        CASE 
            WHEN c.is_group = FALSE THEN (
                SELECT u.username 
                FROM User u 
                INNER JOIN ConversationParticipants cp2 ON u.user_id = cp2.user_id 
                WHERE cp2.conversation_id = c.conversation_id 
                AND cp2.user_id != ?
            )
            ELSE c.name 
        END as display_name,
        c.photo,
        c.is_group,
        c.created_at
    FROM Conversation c
    INNER JOIN ConversationParticipants cp ON c.conversation_id = cp.conversation_id
    WHERE cp.user_id = ?`

	rows, err := db.Query(conversationsQuery, userID, userID)
	if err != nil {
		log.Printf("Error querying conversations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation

	for rows.Next() {
		var conv models.Conversation
		var name string

		err := rows.Scan(
			&conv.ID,
			&name,
			&conv.Photo,
			&conv.Is_group,
			&conv.Created_at,
		)
		if err != nil {
			return nil, err
		}

		conv.Name = name

		// Query for participants
		participantsQuery := `
            SELECT 
                u.user_id,
                u.username,
                u.icon
            FROM ConversationParticipants cp
            JOIN User u ON cp.user_id = u.user_id
            WHERE cp.conversation_id = ?`

		participantRows, err := db.Query(participantsQuery, conv.ID)
		if err != nil {
			return nil, err
		}
		defer participantRows.Close()

		var participants []models.User
		for participantRows.Next() {
			var user models.User
			err := participantRows.Scan(
				&user.ID,
				&user.Username,
				&user.Icon,
			)
			if err != nil {
				return nil, err
			}
			participants = append(participants, user)
		}
		conv.Participants = participants

		// Query for messages with replied to messages
		messagesQuery := `
            SELECT 
                m.message_id,
                m.content,
                m.content_type,
                m.sent_time,
                m.edited_time,
                m.deleted_time,
                m.replied_to,
                m.is_forwarded,
                u.user_id,
                u.username,
                u.icon,
                rm.message_id,
                rm.content,
                rm.content_type
            FROM Message m
            JOIN User u ON m.sender_id = u.user_id
            LEFT JOIN Message rm ON m.replied_to = rm.message_id
            WHERE m.conversation_id = ?
            ORDER BY m.sent_time ASC`

		messageRows, err := db.Query(messagesQuery, conv.ID)
		if err != nil {
			log.Printf("Error querying messages: %v", err)
			continue
		}
		defer messageRows.Close()

		var messages []models.Message
		for messageRows.Next() {
			var msg models.Message
			var sender models.User
			var content []byte
			var editedTime, deletedTime sql.NullTime
			var repliedTo sql.NullInt64
			var isForwarded bool
			var repliedToMsg struct {
				ID          sql.NullInt64
				Content     []byte
				ContentType sql.NullString
			}

			err := messageRows.Scan(
				&msg.ID,
				&content,
				&msg.ContentType,
				&msg.SentTime,
				&editedTime,
				&deletedTime,
				&repliedTo,
				&isForwarded,
				&sender.ID,
				&sender.Username,
				&sender.Icon,
				&repliedToMsg.ID,
				&repliedToMsg.Content,
				&repliedToMsg.ContentType,
			)
			if err != nil {
				log.Printf("Error scanning message: %v", err)
				continue
			}

			msg.Content = content
			msg.EditedTime = editedTime
			msg.DeletedTime = deletedTime
			msg.RepliedTo = repliedTo
			msg.IsForwarded = isForwarded
			msg.ConversationID = conv.ID
			msg.Sender = sender

			if repliedToMsg.ID.Valid {
				msg.RepliedToMessage = &models.Message{
					ID:          int(repliedToMsg.ID.Int64),
					Content:     repliedToMsg.Content,
					ContentType: repliedToMsg.ContentType.String,
				}
			}

			messages = append(messages, msg)
		}

		// Query for reactions
		reactionQuery := `
            SELECT 
                r.reaction_id,
                r.message_id,
                r.user_id,
                r.reaction,
                u.user_id,
                u.username,
                u.icon
            FROM Reactions r
            JOIN User u ON r.user_id = u.user_id
            WHERE r.message_id IN (
                SELECT message_id 
                FROM Message 
                WHERE conversation_id = ?
            )`

		reactionRows, err := db.Query(reactionQuery, conv.ID)
		if err != nil {
			log.Printf("Error querying reactions: %v", err)
			conv.Messages = messages
			conversations = append(conversations, conv)
			continue
		}
		defer reactionRows.Close()

		reactionMap := make(map[int][]models.Reaction)
		for reactionRows.Next() {
			var r models.Reaction
			var u models.User
			err := reactionRows.Scan(
				&r.ID,
				&r.MessageID,
				&r.UserID,
				&r.Reaction,
				&u.ID,
				&u.Username,
				&u.Icon,
			)
			if err != nil {
				log.Printf("Error scanning reaction: %v", err)
				continue
			}
			r.User = u
			reactionMap[r.MessageID] = append(reactionMap[r.MessageID], r)
		}

		// Attach reactions to messages
		for i := range messages {
			if reactions, ok := reactionMap[messages[i].ID]; ok {
				messages[i].Reactions = reactions
			}
		}

		conv.Messages = messages
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}
