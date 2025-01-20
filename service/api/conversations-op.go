package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wasa/service/shared/helper"
	"wasa/service/shared/models"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type reqMessageBody struct {
	Content                   string `json:"content"`
	ContentType               string `json:"content_type"`
	RepliedTo                 *int   `json:"replied_to,omitempty"`
	IsForwarded               bool   `json:"is_forwarded,omitempty"`
	DestinationConversationID *int   `json:"destination_conversation_id,omitempty"`
}

// func (rt *APIRouter) a(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// }

func (rt *APIRouter) conversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	conversations, err := rt.db.GetAllConversations(user_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(conversations)
}

func (rt *APIRouter) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	messages, err := rt.db.GetMessagesFromConversation(conversationID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	if messages == nil {
		messages = []models.Message{}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(messages)
}

func (rt *APIRouter) deleteConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))

	count, _ := rt.db.CountParticipants(conversationID)
	if count <= 2 {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	_, err = rt.db.RemoveUserFromConversation(conversationID, user_id)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode("Could not remove user from conversation")
		return
	}

	if count == 1 {
		_, err = rt.db.DeleteConversation(conversationID)
		if err != nil {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode("Could not delete conversation.")
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode("succesfully left group")

}

func (rt *APIRouter) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))
	message_id, _ := strconv.Atoi(ps.ByName("message_id"))

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	exists, _ = rt.db.IsMessageFromUser(message_id, user_id)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Message not from this user.")
		return
	}

	deletedMessage, err := rt.db.DeleteMessage(message_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Couldn't delete message.")
		return
	}

	wsMessage := WebSocketMessage{
		Type:           "message_deletion",
		ConversationID: deletedMessage.ConversationID,
		Payload: map[string]interface{}{
			"message": deletedMessage,
		},
		Timestamp: time.Now(),
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err == nil {
		rt.baseLogger.WithFields(logrus.Fields{
			"conversation_id": deletedMessage.ConversationID,
			"message_type":    "new_message",
			"recipient_count": len(rt.wsHub.conversationClients[deletedMessage.ConversationID]),
		}).Debug("Broadcasting message via WebSocket")
		rt.wsHub.SendToConversation(deletedMessage.ConversationID, messageJSON)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(deletedMessage)

}

func (rt *APIRouter) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))
	message_id, _ := strconv.Atoi(ps.ByName("message_id"))

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	var req struct {
		Reaction string `json:"reaction"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	new_reaction, err := rt.db.CommentMessage(user_id, message_id, req.Reaction)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	wsMessage := WebSocketMessage{
		Type:           "reaction_update",
		ConversationID: conversationID,
		Payload: map[string]interface{}{
			"reaction": new_reaction,
		},
		Timestamp: time.Now(),
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err == nil {
		rt.baseLogger.WithFields(logrus.Fields{
			"conversation_id": conversationID,
			"message_type":    "reaction_update",
			"recipient_count": len(rt.wsHub.conversationClients[conversationID]),
		}).Debug("Broadcasting reaction via WebSocket")
		rt.wsHub.SendToConversation(conversationID, messageJSON)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(new_reaction)

}

func (rt *APIRouter) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))
	reaction_id, _ := strconv.Atoi(ps.ByName("reaction_id"))

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	_, err = rt.db.IsReactionFromUser(user_id, reaction_id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Reaction is not from the user.")
		return
	}

	deletedReaction, err := rt.db.UncommentMessage(reaction_id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	wsMessage := WebSocketMessage{
		Type:           "reaction_deletion",
		ConversationID: conversationID,
		Payload: map[string]interface{}{
			"reaction": deletedReaction,
		},
		Timestamp: time.Now(),
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err == nil {
		rt.baseLogger.WithFields(logrus.Fields{
			"conversation_id": conversationID,
			"message_type":    "reaction_deletion",
			"recipient_count": len(rt.wsHub.conversationClients[conversationID]),
		}).Debug("Broadcasting reaction via WebSocket")
		rt.wsHub.SendToConversation(conversationID, messageJSON)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode("succesfully uncommented message")

}

func (rt *APIRouter) updateGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	conversation_id, _ := strconv.Atoi(ps.ByName("conversation_id"))

	isGroup, err := rt.db.IsGroup(conversation_id)
	if err != nil || !isGroup {
		http.Error(w, "Not a group conversation", http.StatusBadRequest)
		return
	}

	exists, _ := rt.db.IsUserInConversation(user_id, conversation_id)
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_, err = rt.db.UpdateGroupName(conversation_id, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Group name updated successfully")
}

func (rt *APIRouter) UpdateGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	conversation_id, _ := strconv.Atoi(ps.ByName("conversation_id"))

	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := rt.db.UpdateGroupPhoto(conversation_id, imageData)
	if !result {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Icon not accepted.")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Error")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Icon updated succesfully")
}

func (rt *APIRouter) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var req struct {
		Members []int  `json:"members"`
		Name    string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Members) < 1 {
		http.Error(w, "Conversation must have at least one other member", http.StatusBadRequest)
		return
	}

	conversation, err := rt.db.CreateConversation(req.Members, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, memberID := range req.Members {
		if clients, ok := rt.wsHub.userConnections[memberID]; ok {
			for _, client := range clients {
				rt.wsHub.AddConversationClient(conversation.ID, client)
			}
		}
	}

	wsMessage := WebSocketMessage{
		Type:           "new_conversation",
		ConversationID: conversation.ID,
		Payload: map[string]interface{}{
			"conversation": conversation,
		},
		Timestamp: time.Now(),
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err == nil {
		for _, memberID := range req.Members {
			rt.wsHub.SendToUser(memberID, messageJSON)
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(conversation)
}

func (rt *APIRouter) addGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversation_id, _ := strconv.Atoi(ps.ByName("conversation_id"))

	isGroup, err := rt.db.IsGroup(conversation_id)
	if err != nil || !isGroup {
		http.Error(w, "Not a group conversation", http.StatusBadRequest)
		return
	}

	exists, _ := rt.db.IsUserInConversation(user_id, conversation_id)
	if !exists {
		http.Error(w, "User not in group", http.StatusUnauthorized)
		return
	}

	var req struct {
		Members []int `json:"members"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := rt.db.AddGroupMembers(conversation_id, req.Members); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conversation, err := rt.db.GetConversation(conversation_id)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to get conversation after adding members")
		http.Error(w, "Failed to get updated conversation", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.WithField("conversation", conversation).Debug("Got conversation data")

	wsMessage := WebSocketMessage{
		Type:           "new_conversation",
		ConversationID: conversation_id,
		Payload: map[string]interface{}{
			"conversation": conversation,
		},
		Timestamp: time.Now(),
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to marshal WebSocket message")
		http.Error(w, "Failed to create WebSocket message", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.WithField("message", string(messageJSON)).Debug("Sending WebSocket message")

	for _, memberID := range req.Members {
		rt.baseLogger.WithFields(logrus.Fields{
			"member_id":       memberID,
			"conversation_id": conversation_id,
		}).Debug("Sending to user")

		rt.wsHub.SendToUser(memberID, messageJSON)

		if clients, ok := rt.wsHub.userConnections[memberID]; ok {
			for _, client := range clients {
				rt.wsHub.AddConversationClient(conversation_id, client)
				rt.baseLogger.Debug("Added user to conversation clients")
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Members added successfully")
}

func (rt *APIRouter) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	sourceConversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))

	var message models.Message
	message.Sender.ID = userID

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		imageData, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if destID := r.URL.Query().Get("destination_conversation_id"); destID != "" {
			destConvID, err := strconv.Atoi(destID)
			if err != nil {
				http.Error(w, "Invalid destination conversation ID", http.StatusBadRequest)
				return
			}
			message.ConversationID = destConvID
			message.IsForwarded = true
		} else {
			message.ConversationID = sourceConversationID
		}

		message.Content = imageData
		message.ContentType = "image"
	} else {
		var reqBody reqMessageBody
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		defer r.Body.Close()

		if reqBody.ContentType == "" {
			reqBody.ContentType = "text"
		}

		if reqBody.DestinationConversationID != nil {
			message.ConversationID = *reqBody.DestinationConversationID
			message.IsForwarded = true
			isInSource, err := rt.db.IsUserInConversation(userID, sourceConversationID)
			if err != nil || !isInSource {
				http.Error(w, "Not authorized to access source conversation", http.StatusForbidden)
				return
			}
		} else {
			message.ConversationID = sourceConversationID
		}

		message.Content = []byte(reqBody.Content)
		message.ContentType = reqBody.ContentType
		message.IsForwarded = reqBody.IsForwarded

		if reqBody.RepliedTo != nil {
			message.RepliedTo = helper.PtrToNullInt64(reqBody.RepliedTo)
		}
	}

	exists, _ := rt.db.ConversationExists(message.ConversationID)
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Conversation does not exist.")
		return
	}

	isInDest, err := rt.db.IsUserInConversation(userID, message.ConversationID)
	if err != nil || !isInDest {
		http.Error(w, "Not authorized to send to this conversation", http.StatusForbidden)
		return
	}

	var insertedMessage interface{}
	if message.RepliedTo.Valid {
		insertedMessage, err = rt.db.ReplyToMessage(message)
	} else if message.IsForwarded {
		insertedMessage, err = rt.db.ForwardMessage(message)
	} else {
		insertedMessage, err = rt.db.SendMessage(message)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wsMessage := WebSocketMessage{
		Type:           "new_message",
		ConversationID: message.ConversationID,
		Payload: map[string]interface{}{
			"message": insertedMessage,
		},
		Timestamp: time.Now(),
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err == nil {
		rt.baseLogger.WithFields(logrus.Fields{
			"conversation_id": message.ConversationID,
			"message_type":    "new_message",
			"recipient_count": len(rt.wsHub.conversationClients[message.ConversationID]),
		}).Debug("Broadcasting message via WebSocket")
		rt.wsHub.SendToConversation(message.ConversationID, messageJSON)
	}

	rt.baseLogger.WithFields(logrus.Fields{
		"conversation_id": message.ConversationID,
		"message_type":    "new_message",
	}).Debug("Broadcasting message via WebSocket")

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(insertedMessage)
}
