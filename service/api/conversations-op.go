package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wasa/service/api/responses"
	"wasa/service/shared/models"

	"github.com/julienschmidt/httprouter"
)

type reqMessageBody struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	RepliedTo   *int   `json:"replied_to"`
}

// func (rt *APIRouter) a(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// }

func (rt *APIRouter) conversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := getToken(r)
	conversations := rt.db.GetAllConversations(id)

	if conversations == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Problem getting conversations")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func (rt *APIRouter) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))

	exists, _ := rt.db.IsUserInConversation(getToken(r), conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	messages := rt.db.GetMessagesFromConversation(conversationID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (rt *APIRouter) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var message models.Message
	message.SenderID = getToken(r)
	message.ConversationID, _ = strconv.Atoi(ps.ByName("conversation_id"))

	var reqBody reqMessageBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer r.Body.Close()

	if reqBody.ContentType == "" {
		reqBody.ContentType = "text"
	}

	message.Content = []byte(reqBody.Content)
	message.ContentType = reqBody.ContentType

	if reqBody.RepliedTo != nil {
		message.RepliedTo = sql.NullInt64{
			Int64: int64(*reqBody.RepliedTo),
			Valid: true,
		}
	} else {
		message.RepliedTo.Int64 = -1
	}

	id, err := rt.db.SendMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func (rt *APIRouter) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	source_conversation_id, err := strconv.Atoi(ps.ByName("conversation_id"))
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}
	userID := getToken(r)

	var reqBody struct { // ID of message to forward
		Destination_conversation_id int    `json:"destination_conversation_id"` // Where to forward to
		Content                     string `json:"content"`
		Content_type                string `json:"content_type"`
		Original_message_id         int    `json:"original_message_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isInSource, err := rt.db.IsUserInConversation(source_conversation_id, userID)
	if err != nil || !isInSource {
		http.Error(w, "Not authorized to access source conversation", http.StatusForbidden)
		return
	}

	isInDest, err := rt.db.IsUserInConversation(reqBody.Destination_conversation_id, userID)
	if err != nil || !isInDest {
		http.Error(w, "Not authorized to forward to destination conversation", http.StatusForbidden)
		return
	}

	var message models.Message
	message.ConversationID = reqBody.Destination_conversation_id
	message.Content = []byte(reqBody.Content)
	message.ContentType = reqBody.Content_type
	message.SenderID = userID
	message.ForwardedFrom = sql.NullInt64{
		Int64: int64(reqBody.Original_message_id),
		Valid: true,
	}

	id, err := rt.db.SendMessage(message)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(id)

}

func (rt *APIRouter) deleteConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id := getToken(r)
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	count, err := rt.db.CountParticipants(conversationID)

	if count <= 2 {
		responses.SendError(w, errors.New("can't delete a direct conversation"), http.StatusBadRequest)
		return
	}

	_, err = rt.db.RemoveUserFromConversation(conversationID, user_id)
	if err != nil {
		responses.SendError(w, err, http.StatusInternalServerError)
		return
	}

	if count == 1 {
		_, err = rt.db.DeleteConversation(conversationID)
		if err != nil {
			responses.SendError(w, err, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("succesfully left group")

}

func (rt *APIRouter) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id := getToken(r)
	conversationID, _ := strconv.Atoi(ps.ByName("conversation_id"))
	message_id, _ := strconv.Atoi(ps.ByName("message_id"))

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("User is not on the conversation.")
		return
	}

	// todo: might need to check if message is from conversation for API calls
	// since, from the UI, we always will delete from the conversation itself,
	// might not need it

	exists, _ = rt.db.IsMessageFromUser(message_id, user_id)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Message not from this user.")
		return
	}

	_, err := rt.db.DeleteMessage(message_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Couldn't delete message.")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("succesfully deleted message")

}
