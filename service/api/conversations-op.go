package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wasa/service/api/responses"
	"wasa/service/shared/helper"
	"wasa/service/shared/models"

	"github.com/julienschmidt/httprouter"
)

type reqMessageBody struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	RepliedTo   *int   `json:"replied_to,omitempty"`
	IsForwarded *int `json:"isForwarded,omitempty"`
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
	conversations := rt.db.GetAllConversations(user_id)

	if conversations == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Problem getting conversations")
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

	messages := rt.db.GetMessagesFromConversation(conversationID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(messages)
}

func (rt *APIRouter) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	var message models.Message
	message.Sender.ID = user_id
	message.ConversationID, _ = strconv.Atoi(ps.ByName("conversation_id"))
	exists, _ := rt.db.ConversationExists(message.ConversationID)
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Conversation do not exist.")
		return
	}

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

	message.Content = []byte(reqBody.Content)
	message.ContentType = reqBody.ContentType

	insertedMessage, err := rt.db.SendMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(insertedMessage)
}

func (rt *APIRouter) replyToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	var message models.Message
	message.Sender.ID = user_id
	message.ConversationID, _ = strconv.Atoi(ps.ByName("conversation_id"))

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

	message.Content = []byte(reqBody.Content)
	message.ContentType = reqBody.ContentType
	message.RepliedTo = helper.PtrToNullInt64(reqBody.RepliedTo)

	insertedMessage, err := rt.db.ReplyToMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(insertedMessage)
}

func (rt *APIRouter) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	source_conversation_id, err := strconv.Atoi(ps.ByName("conversation_id"))
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}
	user_id, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	var reqBody struct {
		Destination_conversation_id int    `json:"destination_conversation_id"` // Where to forward to
		Content                     string `json:"content"`
		Content_type                string `json:"content_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isInSource, err := rt.db.IsUserInConversation(user_id, source_conversation_id)
	if err != nil || !isInSource {
		http.Error(w, "Not authorized to access source conversation", http.StatusForbidden)
		return
	}

	// destination_conversation_id := reqBody.Destination_conversation_id
	isInDest, err := rt.db.IsUserInConversation(user_id, reqBody.Destination_conversation_id)
	if err != nil || !isInDest {
		http.Error(w, "Not authorized to forward to destination conversation", http.StatusForbidden)
		return
	}

	var message models.Message
	message.ConversationID = reqBody.Destination_conversation_id
	message.Content = []byte(reqBody.Content)
	message.ContentType = reqBody.Content_type
	message.Sender.ID = user_id
	message.IsForwarded = true

	id, err := rt.db.ForwardMessage(message)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(id)

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

	exists, _ := rt.db.IsUserInConversation(user_id, conversationID)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("User is not on the conversation.")
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

	// todo: might need to check if message is from conversation for API calls
	// since, from the UI, we always will delete from the conversation itself,
	// might not need it

	exists, _ = rt.db.IsMessageFromUser(message_id, user_id)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Message not from this user.")
		return
	}

	_, err = rt.db.DeleteMessage(message_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode("Couldn't delete message.")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode("succesfully deleted message")

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

	reactionBytes := []byte(req.Reaction)

	_, err = rt.db.CommentMessage(user_id, message_id, reactionBytes)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode("succesfully commented message")

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

	_, err = rt.db.UncommentMessage(reaction_id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
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
		Photo string `json:"photo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_, err = rt.db.UpdateGroupPhoto(conversation_id, req.Photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Group photo updated successfully")
}

func (rt *APIRouter) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var req struct {
		Members []int `json:"members"`
		Name string `json:"name"`
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


	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(conversation)
}

func (rt *APIRouter) addGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Members added successfully")
}
