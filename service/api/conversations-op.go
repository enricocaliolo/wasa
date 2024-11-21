package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
