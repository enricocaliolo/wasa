package conversationDB

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"wasa/service/shared/models"
)

func GetAllConversations(db *sql.DB, user_id int) []models.Conversation {
    query := `SELECT 
    c.conversation_id,
    COALESCE(c.name, '') as name,
    COALESCE(c.is_group, false) as is_group,
    COALESCE(c.created_at, '') as created_at,
    GROUP_CONCAT(
        CASE 
            WHEN cp2.user_id = cp1.user_id THEN NULL 
            ELSE cp2.user_id 
        END
    ) as participant_ids,
    GROUP_CONCAT(
        CASE 
            WHEN cp2.user_id = cp1.user_id THEN NULL 
            ELSE COALESCE(cp2.joined_at, '')
        END
    ) as joined_ats,
    GROUP_CONCAT(
        CASE 
            WHEN cp2.user_id = cp1.user_id THEN NULL 
            ELSE u.username 
        END
    ) as participant_names
    FROM ConversationParticipants cp1
    JOIN Conversation c ON cp1.conversation_id = c.conversation_id
    JOIN ConversationParticipants cp2 ON c.conversation_id = cp2.conversation_id
    JOIN User u ON cp2.user_id = u.user_id
    WHERE cp1.user_id = ?
    GROUP BY c.conversation_id;`


    rows, err := db.Query(query, user_id)
    if err != nil {
        log.Fatal(err)
        return nil
    }
    defer rows.Close()

    var conversations []models.Conversation

    for rows.Next() {
        var conversation models.Conversation
        var participantIdsStr, joinedAtsStr, participantNamesStr string

        err := rows.Scan(
            &conversation.ID,
            &conversation.Name,
            &conversation.Is_group,
            &conversation.Created_at,
            &participantIdsStr,
            &joinedAtsStr,
            &participantNamesStr,  // Add this line
        )
        if err != nil {
            log.Fatal(err)
            return nil
        }

        // Convert comma-separated strings to slices
        participantIds := strings.Split(participantIdsStr, ",")
        joinedAts := strings.Split(joinedAtsStr, ",")
        participantNames := strings.Split(participantNamesStr, ",")  // Add this line
        
        // Create participants slice
        conversation.ConversationParticipant = make([]models.ConversationParticipant, 0, len(participantIds))
        
        for i, idStr := range participantIds {
            id, err := strconv.Atoi(idStr)
            if err != nil {
                log.Fatal(err)
                return nil
            }
            
            participant := models.ConversationParticipant{
                User_id:   id,
                Joined_at: joinedAts[i],
                Name:      participantNames[i],
            }
            conversation.ConversationParticipant = append(conversation.ConversationParticipant, participant)
        }

        conversations = append(conversations, conversation)
    }

    if err = rows.Err(); err != nil {
        log.Fatal(err)
        return nil
    }

    return conversations
}