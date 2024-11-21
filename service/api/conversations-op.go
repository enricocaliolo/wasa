package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"wasa/service/shared/models"

	"github.com/julienschmidt/httprouter"
)

type reqMessageBody struct {
	Content       string `json:"content"`
	ContentType   string `json:"content_type"`
	RepliedTo     *int   `json:"replied_to"`
	ForwardedFrom *int   `json:"forwarded_from"`
}

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
	conversationID, _ := strconv.Atoi(ps.ByName("id"))

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
	message.ConversationID, _ = strconv.Atoi(ps.ByName("id"))

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

	if reqBody.ForwardedFrom != nil {
		message.ForwardedFrom = sql.NullInt64{
			Int64: int64(*reqBody.ForwardedFrom),
			Valid: true,
		}
	} else {
		message.ForwardedFrom.Int64 = -1
	}

	id, err := rt.db.SendMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(id)
}
